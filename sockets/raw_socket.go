package sockets

import (
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/NamelessOne91/bisturi/protocols"
)

const mask = 0xff00

// hostToNetworkShort converts a short (uint16) from host (usually Little Endian)
// to network (Big Endian) byte order
func hostToNetworkShort(i uint16) uint16 {
	return (i<<8)&mask | i>>8
}

type NetworkPacket interface {
	Source() string
	Destination() string
	Info() string
}

// RawSocket represents a raw socket and stores info about its file descriptor,
// Ethernet protocol type and Link Layer info
type RawSocket struct {
	shutdownChan chan os.Signal
	fd           int
	ethType      uint16
	layer4Filter string
	sll          syscall.SockaddrLinklayer
}

// NewRawSocket opens a raw socket for the specified protocol by calling SYS_SOCKET
// and returns the struct representing it, or eventual errors
func NewRawSocket(protocol string, ethType uint16) (*RawSocket, error) {
	filter := "all"
	if strings.HasPrefix(protocol, "udp") {
		filter = "udp"
	} else if strings.HasPrefix(protocol, "tcp") {
		filter = "tcp"
	}

	rawSocket := &RawSocket{
		shutdownChan: make(chan os.Signal, 1),
		ethType:      ethType,
		layer4Filter: filter,
	}
	// AF_PACKET specifies a packet socket, operating at the data link layer (Layer 2)
	// SOCK_RAW specifies a raw socket
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(hostToNetworkShort(ethType)))
	if err != nil {
		return nil, err
	}
	rawSocket.fd = fd

	return rawSocket, nil
}

// BindSocket binds a raw socket  to a network interface allowing to monitor
// and analyze packets traversing it
func (rs *RawSocket) Bind(iface net.Interface) error {
	// network stack uses Big Endian
	rs.sll.Protocol = hostToNetworkShort(rs.ethType)
	rs.sll.Ifindex = iface.Index

	return syscall.Bind(rs.fd, &rs.sll)
}

// ReadToChan calls SYS_RECVFROM to read data traversing the binded network interface and sends its representation to the passed channel.
// Errors are sent to another passed channel
func (rs *RawSocket) ReadToChan(dataChan chan<- NetworkPacket, errChan chan<- error) {
	buf := make([]byte, 4096)

	for {
		n, _, err := syscall.Recvfrom(rs.fd, buf, 0)
		if err != nil {
			errChan <- fmt.Errorf("error reading from raw socket: %v", err)
			continue
		}

		switch rs.ethType {
		case syscall.ETH_P_ARP:
			// TODO: ARP parsing
		case syscall.ETH_P_IP, syscall.ETH_P_IPV6:
			packet, err := protocols.IPPacketFromBytes(buf[:n])
			if err != nil {
				errChan <- fmt.Errorf("error reading IP packet: %v", err)
				continue
			}
			// IPv4 VS IPv6 packets filtering should be handled by the socket itself
			l4Protocol := packet.Header().TransportLayerProtocol()
			if rs.layer4Filter == "all" || (l4Protocol == rs.layer4Filter) {
				handleLayer4Protocol(l4Protocol, packet, dataChan, errChan)
			}
		}
	}
}

// Close closes the raw socket by calling SYS_CLOSE on its file descriptor
func (rs *RawSocket) Close() error {
	return syscall.Close(rs.fd)
}

func handleLayer4Protocol(protocol string, packet protocols.IPPacket, dataChan chan<- NetworkPacket, errChan chan<- error) {
	var np NetworkPacket
	var err error

	switch protocol {
	case "udp":
		np, err = protocols.UDPPacketFromIPPacket(packet)
	case "tcp":
		np, err = protocols.TCPPacketFromIPPacket(packet)
	}

	if err != nil {
		errChan <- fmt.Errorf("error reading %s packet: %v", protocol, err)
		return
	}
	dataChan <- np
}
