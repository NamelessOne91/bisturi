package sockets

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/NamelessOne91/bisturi/packets"
)

const mask = 0xff00

// maps the protocol names to Ethernet protocol types values
var protocolEthernetType = map[string]uint16{
	"all":  syscall.ETH_P_ALL,
	"arp":  syscall.ETH_P_ARP,
	"ip":   syscall.ETH_P_IP,
	"ipv6": syscall.ETH_P_IPV6,
	"udp":  syscall.ETH_P_IP, // UDP and TCP are part of IP, need special handling if filtered specifically
	"udp6": syscall.ETH_P_IPV6,
	"tcp":  syscall.ETH_P_IP,
	"tcp6": syscall.ETH_P_IPV6,
}

var errUnsupportedProtocol = errors.New("unsupported protocol")

// hostToNetworkShort converts a short (uint16) from host (usually Little Endian)
// to network (Big Endian) byte order
func hostToNetworkShort(i uint16) uint16 {
	return (i<<8)&mask | i>>8
}

// RawSocket represents a raw socket and stores info about its file descriptor,
// Ethernet protocl type and Link Layer info
type RawSocket struct {
	shutdownChan chan os.Signal
	fd           int
	ethType      uint16
	sll          syscall.SockaddrLinklayer
}

// NewRawSocket opens a raw socket for the specified protocol by calling SYS_SOCKET
// and returns the struct representing it, or eventual errors
func NewRawSocket(protocol string) (*RawSocket, error) {
	ethType, ok := protocolEthernetType[protocol]
	if !ok {
		return nil, errUnsupportedProtocol
	}

	rawSocket := RawSocket{
		shutdownChan: make(chan os.Signal, 1),
		ethType:      ethType,
	}
	// AF_PACKET specifies a packet socket, operating at the data link layer (Layer 2)
	// SOCK_RAW specifies a raw socket
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(hostToNetworkShort(ethType)))
	if err != nil {
		return nil, err
	}
	rawSocket.fd = fd

	return &rawSocket, nil
}

// BindSocket binds a raw socket  to a network interface allowing to monitor
// and analyze packets traversing it
func (rs *RawSocket) Bind(iface net.Interface) error {
	// network stack uses Big Endian
	rs.sll.Protocol = hostToNetworkShort(rs.ethType)
	rs.sll.Ifindex = iface.Index

	// handle graceful shutdown on CTRL + C and similar
	signal.Notify(rs.shutdownChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-rs.shutdownChan
		log.Println("Received interrupt, stopping...")
		if err := syscall.Close(rs.fd); err != nil {
			log.Printf("Failed to close rad socket with file descriptor %d", rs.fd)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	return syscall.Bind(rs.fd, &rs.sll)
}

// ReadPackets calls SYS_RECVFROM to read packets traversing the binded network interface an display info about them.
func (rs *RawSocket) ReadPackets() {
	// read incoming packets
	buf := make([]byte, 4096)
	for {
		n, _, err := syscall.Recvfrom(rs.fd, buf, 0)
		if err != nil {
			log.Println("Error reading from socket:", err)
			continue
		}

		switch rs.ethType {
		case syscall.ETH_P_IP:
			packet, err := packets.IPv4PacketFromBytes(buf[:n])
			if err != nil {
				log.Println("Error reading IPv4 packet:", err)
				continue
			}
			log.Println(packet.Info())
		case syscall.ETH_P_IPV6:
			packet, err := packets.IPv6PacketFromBytes(buf[:n])
			if err != nil {
				log.Println("Error reading IPv6 packet:", err)
				continue
			}
			log.Println(packet.Info())
		}
	}
}

// Close closes the raw socket by calling SYS_CLOSE on its file descriptor
func (rs *RawSocket) Close() error {
	return syscall.Close(rs.fd)
}
