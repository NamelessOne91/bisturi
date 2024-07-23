package protocols

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type TCPPacket struct {
	IPPacket IPPacket
	Header   TCPHeader
}

type TCPHeader struct {
	SourcePort      uint16
	DestinationPort uint16
	SequenceNumber  uint32
	AckNumber       uint32
	RawOffset       uint8 // rawOffset is in 4-byte words
	Flags           uint8
	WindowSize      uint16
	Checksum        uint16
	UrgentPointer   uint16
	Options         []byte
}

var (
	ErrTCPHeaderTooShort    = errors.New("TCP header must be at least 20 bytes")
	ErrTCPHeaderLenMismatch = errors.New("TCP header length less than raw Offset")
)

func TCPPacketFromIPPacket(ip IPPacket) (*TCPPacket, error) {
	TCPHeader, err := TCPHeaderFromBytes(ip.Payload())
	if err != nil {
		return nil, err
	}

	return &TCPPacket{
		IPPacket: ip,
		Header:   *TCPHeader,
	}, nil
}

func TCPHeaderFromBytes(raw []byte) (*TCPHeader, error) {
	if len(raw) < 20 {
		return nil, ErrTCPHeaderTooShort
	}

	offset := raw[12] >> 4
	hLen := int(offset) * 4
	if len(raw) < hLen {
		return nil, ErrTCPHeaderLenMismatch
	}

	h := &TCPHeader{
		SourcePort:      binary.BigEndian.Uint16(raw[0:2]),
		DestinationPort: binary.BigEndian.Uint16(raw[2:4]),
		SequenceNumber:  binary.BigEndian.Uint32(raw[4:8]),
		AckNumber:       binary.BigEndian.Uint32(raw[8:12]),
		RawOffset:       offset,
		Flags:           raw[13],
		WindowSize:      binary.BigEndian.Uint16(raw[14:16]),
		Checksum:        binary.BigEndian.Uint16(raw[16:18]),
		UrgentPointer:   binary.BigEndian.Uint16(raw[18:20]),
	}
	if hLen > 20 {
		h.Options = raw[20:hLen]
	}
	return h, nil
}

// Info return an human-readable string containing the main TCP packet data
func (p TCPPacket) Info() string {
	return fmt.Sprintf("%s - port %d to port %d",
		p.IPPacket.Info(), p.Header.SourcePort, p.Header.DestinationPort,
	)
}

func (p TCPPacket) Source() string {
	return fmt.Sprintf("%s:%d", p.IPPacket.Header().Source(), p.Header.SourcePort)
}

func (p TCPPacket) Destination() string {
	return fmt.Sprintf("%s:%d", p.IPPacket.Header().Destination(), p.Header.DestinationPort)
}
