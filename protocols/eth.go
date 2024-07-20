package protocols

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// maps the ETH frame values to the corresponding string representation
var etherTypesValues = map[uint16]string{
	0x0800: "IPv4",
	0x0806: "ARP",
	0x0842: "Wake-on-LAN",
	0x86DD: "IPv6",
	0x8808: "Ethernet flow control",
}

// EthernetFrame contains Layer 2 data
type EthernetFrame struct {
	destinationMAC net.HardwareAddr
	sourceMAC      net.HardwareAddr
	etherType      uint16
	payload        []byte
}

var errInvalidETHFrame = errors.New("ethernet frame header must be 14 bytes")

// EthFrameFromBytes parses an array of bytes to the corresponding ETH frame and returns a pointer to it.
// Returns an error if the number of bytes is less than 14
func EthFrameFromBytes(raw []byte) (*EthernetFrame, error) {
	if len(raw) < 14 {
		return nil, errInvalidETHFrame
	}

	return &EthernetFrame{
		destinationMAC: net.HardwareAddr(raw[0:6]),
		sourceMAC:      net.HardwareAddr(raw[6:12]),
		etherType:      binary.BigEndian.Uint16(raw[12:14]),
		payload:        raw[14:],
	}, nil
}

func (f EthernetFrame) EtherType() string {
	return etherTypesValues[f.etherType]
}

// Info returns an human-readable string containing all the ETH frame data
func (f EthernetFrame) Info() string {
	return fmt.Sprintf("%s Ethernet Frame from MAC %s to MAC %s", f.EtherType(), f.sourceMAC, f.destinationMAC)
}
