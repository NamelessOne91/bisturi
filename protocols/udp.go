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

func UDPPacketFromIPPacket(ip IPPacket) (*UDPPacket, error) {
	udpHeader, err := udpHeaderFromBytes(ip.Payload())

	return &UDPPacket{
		ipPacket: ip,
		header:   *udpHeader,
	}, err
}

func udpv4PacketFromBytes(raw []byte) (*UDPPacket, error) {
	ipPacket, err := ipv4PacketFromBytes(raw)
	if err != nil {
		return nil, err
	}

	// 14 bytes eth frame + (IPv4 IHL * 4)
	offset := 14 + (ipPacket.header.ihl * 4)
	h, err := udpHeaderFromBytes(raw[offset:])
	if err != nil {
		return nil, err
	}

	return &UDPPacket{
		ipPacket: ipPacket,
		header:   *h,
	}, nil
}

func udpv6PacketFromBytes(raw []byte) (*UDPPacket, error) {
	ipPacket, err := ipv6PacketFromBytes(raw)
	if err != nil {
		return nil, err
	}
	// 14 bytes eth frame + 40 IPv6 header
	h, err := udpHeaderFromBytes(raw[54:])
	if err != nil {
		return nil, err
	}

	return &UDPPacket{
		ipPacket: ipPacket,
		header:   *h,
	}, nil
}

func (p UDPPacket) Info() string {
	return fmt.Sprintf(`
UDP packet

Source Port: %d
Destination Port: %d
Length: %d
Checksum: %d

===============================
%s
===============================
`,
		p.header.sourcePort, p.header.sourcePort, p.header.length, p.header.checksum, p.ipPacket.Info(),
	)
}

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