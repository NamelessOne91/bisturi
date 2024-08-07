package protocols

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// maps the ETH frame values to the corresponding string representation
var EtherTypesValues = map[uint16]string{
	0x0800: "IPv4",
	0x0806: "ARP",
	0x0842: "Wake-on-LAN",
	0x8535: "RARP",
	0x86DD: "IPv6",
	0x8808: "Ethernet flow control",
}

// EthernetFrame contains Layer 2 data
type EthernetFrame struct {
	DestinationMAC net.HardwareAddr
	SourceMAC      net.HardwareAddr
	EtherType      uint16
	Payload        []byte
}

var errInvalidETHFrame = errors.New("ethernet frame header must be 14 bytes")

// EthFrameFromBytes parses an array of bytes to the corresponding ETH frame and returns a pointer to it.
// Returns an error if the number of bytes is less than 14
func EthFrameFromBytes(raw []byte) (*EthernetFrame, error) {
	if len(raw) < 14 {
		return nil, errInvalidETHFrame
	}

	return &EthernetFrame{
		DestinationMAC: net.HardwareAddr(raw[0:6]),
		SourceMAC:      net.HardwareAddr(raw[6:12]),
		EtherType:      binary.BigEndian.Uint16(raw[12:14]),
		Payload:        raw[14:],
	}, nil
}

func (f EthernetFrame) Type() string {
	return EtherTypesValues[f.EtherType]
}

// Info returns an human-readable string containing all the ETH frame data
func (f EthernetFrame) Info() string {
	etv := EtherTypesValues[f.EtherType]

	return fmt.Sprintf(`
Ethernet Frame

Destination MAC: %s
Source MAC: %s
EtherType: 0x%X (%s)`,
		f.DestinationMAC, f.SourceMAC, f.EtherType, etv,
	)
}
