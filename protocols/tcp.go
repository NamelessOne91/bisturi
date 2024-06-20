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
	rawOffset       uint8
	reserved        uint8
	flags           uint16
	windowSize      uint16
	checksum        uint16
	urgentPointer   uint16
	options         []byte
}

var errTCPHeaderTooShort = errors.New("TCP header must be at least 20 bytes")
var errTCPHeaderLenMismatch = errors.New("TCP header length less than raw Offset")

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
	offset := (raw[12] >> 4)
	hLen := int(offset * 4)
	if len(raw) < hLen {
		return nil, errTCPHeaderLenMismatch
	}
	return &tcpHeader{
		sourcePort:      binary.BigEndian.Uint16(raw[0:2]),
		destinationPort: binary.BigEndian.Uint16(raw[2:4]),
		sequenceNumber:  binary.BigEndian.Uint32(raw[4:8]),
		ackNumber:       binary.BigEndian.Uint32(raw[8:12]),
		rawOffset:       offset,
		reserved:        (raw[12] & 0x0E) >> 1,
		flags:           binary.BigEndian.Uint16(raw[12:14]) & 0x1FF,
		windowSize:      binary.BigEndian.Uint16(raw[14:16]),
		checksum:        binary.BigEndian.Uint16(raw[16:18]),
		urgentPointer:   binary.BigEndian.Uint16(raw[18:20]),
		options:         raw[20:hLen],
	}, nil
}

func (p TCPPacket) Info() string {
	return fmt.Sprintf(`
TCP packet

Source Port: %d
Destination Port: %d
Checksum: %d

===============================
%s
===============================
`,
		p.header.sourcePort, p.header.sourcePort, p.header.checksum, p.ipPacket.Info(),
	)
}
