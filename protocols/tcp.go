package protocols

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type TCPPacket struct {
	ipPacket IPPacket
	header   tcpHeader
}

type tcpHeader struct {
	sourcePort      uint16
	destinationPort uint16
	sequenceNumber  uint32
	ackNumber       uint32
	rawOffset       uint8 // rawOffset is in 4-byte words
	flags           uint8
	windowSize      uint16
	checksum        uint16
	urgentPointer   uint16
	options         []byte
}

var (
	errTCPHeaderTooShort    = errors.New("TCP header must be at least 20 bytes")
	errTCPHeaderLenMismatch = errors.New("TCP header length less than raw Offset")
)

func TCPPacketFromIPPacket(ip IPPacket) (*TCPPacket, error) {
	tcpHeader, err := tcpHeaderFromBytes(ip.Payload())

	return &TCPPacket{
		ipPacket: ip,
		header:   *tcpHeader,
	}, err
}

func tcpHeaderFromBytes(raw []byte) (*tcpHeader, error) {
	if len(raw) < 20 {
		return nil, errTCPHeaderTooShort
	}

	offset := raw[12] >> 4
	hLen := int(offset) * 4
	if len(raw) < hLen {
		return nil, errTCPHeaderLenMismatch
	}

	return &tcpHeader{
		sourcePort:      binary.BigEndian.Uint16(raw[0:2]),
		destinationPort: binary.BigEndian.Uint16(raw[2:4]),
		sequenceNumber:  binary.BigEndian.Uint32(raw[4:8]),
		ackNumber:       binary.BigEndian.Uint32(raw[8:12]),
		rawOffset:       offset,
		flags:           raw[13],
		windowSize:      binary.BigEndian.Uint16(raw[14:16]),
		checksum:        binary.BigEndian.Uint16(raw[16:18]),
		urgentPointer:   binary.BigEndian.Uint16(raw[18:20]),
		options:         raw[20:hLen],
	}, nil
}

// Info return an human-readable string containing the main TCP packet data
func (p TCPPacket) Info() string {
	return fmt.Sprintf("%s - port %d to port %d",
		p.ipPacket.Info(), p.header.sourcePort, p.header.destinationPort,
	)
}

func (p TCPPacket) Source() string {
	return fmt.Sprintf("%s:%d", p.ipPacket.Header().Source(), p.header.sourcePort)
}

func (p TCPPacket) Destination() string {
	return fmt.Sprintf("%s:%d", p.ipPacket.Header().Destination(), p.header.destinationPort)
}
