package protocols

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var HardwareTypeValues = map[uint16]string{
	1:  "Ethernet",
	6:  "IEE 802 (Token Ring)",
	11: "ATM",
	12: "HDLC",
}

type ARPPacket struct {
	EthFrame        EthernetFrame
	HardwareType    uint16
	ProtocolType    uint16
	HardwareAddrLen uint8
	ProtocolAddrLen uint8
	Operation       uint16
	SenderHWAddr    net.HardwareAddr
	SenderProtoAddr net.IP // only IPv4 - IPv6 should rely on NDP
	TargetHWAddr    net.HardwareAddr
	TargetProtoAddr net.IP // only IPv4 - IPv6 should rely on NDP
}

var errInvalidARPPacket = errors.New("ARP packet must be 28 bytes")

// ARPPacketFromBytes parses an array of bytes to the corresponding ARP packet and returns a pointer to it.
// Returns an error if the number of bytes is less than 28
func ARPPacketFromBytes(raw []byte) (*ARPPacket, error) {
	frame, err := EthFrameFromBytes(raw)
	if err != nil {
		return nil, err
	}

	if len(frame.Payload) < 28 {
		return nil, errInvalidARPPacket
	}
	payload := frame.Payload

	return &ARPPacket{
		EthFrame:        *frame,
		HardwareType:    binary.BigEndian.Uint16(payload[0:2]),
		ProtocolType:    binary.BigEndian.Uint16(payload[2:4]),
		HardwareAddrLen: payload[4],
		ProtocolAddrLen: payload[5],
		Operation:       binary.BigEndian.Uint16(payload[6:8]),
		SenderHWAddr:    payload[8:14],
		SenderProtoAddr: payload[14:18],
		TargetHWAddr:    payload[18:24],
		TargetProtoAddr: payload[24:28],
	}, nil
}

func (p ARPPacket) Destination() string {
	return fmt.Sprintf("%s|%s", p.TargetHWAddr.String(), p.TargetProtoAddr.String())
}

func (p ARPPacket) Source() string {
	return fmt.Sprintf("%s|%s", p.SenderHWAddr.String(), p.SenderProtoAddr.String())
}

func (p ARPPacket) Info() string {
	return fmt.Sprintf("%s ARP packet from %s to %s", HardwareTypeValues[p.HardwareType], p.Source(), p.Destination())
}
