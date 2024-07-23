package protocols

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type UDPPacket struct {
	IPPacket IPPacket
	Header   UDPHeader
}

type UDPHeader struct {
	SourcePort      uint16
	DestinationPort uint16
	Length          uint16
	Checksum        uint16
}

var errInvalidUDPHeader = errors.New("UDP header must be 8 bytes")

// UDPPacketFromIPPacket parses the passed IPv4 or IPv6 packet's data returning a struct conatining the encapsulated UDP's packet data.
// An error is returned if the headers' constraints are not respected.
func UDPPacketFromIPPacket(ip IPPacket) (*UDPPacket, error) {
	UDPHeader, err := UDPHeaderFromBytes(ip.Payload())
	if err != nil {
		return nil, err
	}

	return &UDPPacket{
		IPPacket: ip,
		Header:   *UDPHeader,
	}, nil
}

// Info return an human-readable string containing the main UDP packet data
func (p UDPPacket) Info() string {
	return fmt.Sprintf("%s - port %d to port %d",
		p.IPPacket.Info(), p.Header.SourcePort, p.Header.DestinationPort,
	)
}

func (p UDPPacket) Source() string {
	return fmt.Sprintf("%s:%d", p.IPPacket.Header().Source(), p.Header.SourcePort)
}

func (p UDPPacket) Destination() string {
	return fmt.Sprintf("%s:%d", p.IPPacket.Header().Destination(), p.Header.DestinationPort)
}

// UDPHeaderFromBytes parses the passed bytes to a struct containing the UDP header data and returns a pointer to it.
// It expects an array of at least 8 bytes
func UDPHeaderFromBytes(raw []byte) (*UDPHeader, error) {
	if len(raw) < 8 {
		return nil, errInvalidUDPHeader
	}

	return &UDPHeader{
		SourcePort:      binary.BigEndian.Uint16(raw[0:2]),
		DestinationPort: binary.BigEndian.Uint16(raw[2:4]),
		Length:          binary.BigEndian.Uint16(raw[4:6]),
		Checksum:        binary.BigEndian.Uint16(raw[6:8]),
	}, nil
}
