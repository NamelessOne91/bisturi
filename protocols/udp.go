package protocols

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type UDPPacket struct {
	ipPacket IPPacket
	header   udpHeader
}

type udpHeader struct {
	sourcePort      uint16
	destinationPort uint16
	length          uint16
	checksum        uint16
}

var errInvalidUDPHeader = errors.New("UDP header must be 8 bytes")

// UDPPacketFromIPPacket parses the passed IPv4 or IPv6 packet's data returning a struct conatining the encapsulated UDP's packet data.
// An error is returned if the headers' constraints are not respected.
func UDPPacketFromIPPacket(ip IPPacket) (*UDPPacket, error) {
	udpHeader, err := udpHeaderFromBytes(ip.Payload())
	if err != nil {
		return nil, err
	}

	return &UDPPacket{
		ipPacket: ip,
		header:   *udpHeader,
	}, nil
}

// Info return an human-readable string containing the main UDP packet data
func (p UDPPacket) Info() string {
	return fmt.Sprintf("%s - port %d to port %d",
		p.ipPacket.Info(), p.header.sourcePort, p.header.destinationPort,
	)
}

func (p UDPPacket) Source() string {
	return fmt.Sprintf("%s:%d", p.ipPacket.Header().Source(), p.header.sourcePort)
}

func (p UDPPacket) Destination() string {
	return fmt.Sprintf("%s:%d", p.ipPacket.Header().Destination(), p.header.destinationPort)
}

// udpHeaderFromBytes parses the passed bytes to a struct containing the UDP header data and returns a pointer to it.
// It expects an array of at least 8 bytes
func udpHeaderFromBytes(raw []byte) (*udpHeader, error) {
	if len(raw) < 8 {
		return nil, errInvalidUDPHeader
	}

	return &udpHeader{
		sourcePort:      binary.BigEndian.Uint16(raw[0:2]),
		destinationPort: binary.BigEndian.Uint16(raw[2:4]),
		length:          binary.BigEndian.Uint16(raw[4:6]),
		checksum:        binary.BigEndian.Uint16(raw[6:8]),
	}, nil
}
