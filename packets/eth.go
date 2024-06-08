package packets

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var etherTypesValues = map[uint16]string{
	0x0800: "IPv4",
	0x0806: "ARP",
	0x0842: "Wake-on-LAN",
	0x86DD: "IPv6",
	0x8808: "Ethernet flow control",
}

type EthernetFrame struct {
	destinationMAC net.HardwareAddr
	sourceMAC      net.HardwareAddr
	etherType      uint16
}

var errInvalidETHFrame = errors.New("ethernet frame header must be 14 bytes")

func EthFrameFromBytes(raw []byte) (*EthernetFrame, error) {
	if len(raw) != 14 {
		return nil, errInvalidETHFrame
	}

	return &EthernetFrame{
		destinationMAC: net.HardwareAddr(raw[0:6]),
		sourceMAC:      net.HardwareAddr(raw[6:12]),
		etherType:      binary.BigEndian.Uint16(raw[12:14]),
	}, nil
}

func (f *EthernetFrame) info() string {
	etv := etherTypesValues[f.etherType]

	return fmt.Sprintf(`Ethernet Frame

Destination MAC: %s
Source MAC: %s
EtherType: 0x%X (%s)`,
		f.destinationMAC, f.sourceMAC, f.etherType, etv,
	)
}
