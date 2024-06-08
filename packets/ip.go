package packets

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var protocolValues = map[uint8]string{
	1:  "ICMP",
	6:  "TCP",
	17: "UDP",
	41: "IPv6",
	58: "ICMPv6",
	89: "OSPF",
}

type IPv4Packet struct {
	ethFrame EthernetFrame
	ipHeader IPv4Header
}

type IPv6Packet struct {
	ethFrame EthernetFrame
	ipHeader IPv6Header
}

type IPv4Header struct {
	version        uint8
	ihl            uint8
	dscp           uint8
	ecn            uint8
	totalLength    uint16
	identification uint16
	flags          uint16
	fragmentOffset uint16
	ttl            uint8
	protocol       uint8
	headerChecksum uint16
	sourceIP       net.IP
	destinationIP  net.IP
	options        []byte
}

type IPv6Header struct {
	version       uint8
	trafficClass  uint8
	flowLabel     uint32
	payloadLength uint16
	nextHeader    uint8
	hopLimit      uint8
	sourceIP      net.IP
	destinationIP net.IP
}

var errIPv4HeaderTooShort = errors.New("IPv4 header must be at least 20 bytes")
var errIPv4HeaderLenLessThanIHL = errors.New("IPv4 header length less than indicated IHL")
var errInvalidIPv6Header = errors.New("IPv6 header must be 40 bytes")

func IPv4PacketFromBytes(raw []byte) (*IPv4Packet, error) {
	frame, err := EthFrameFromBytes(raw[:14])
	if err != nil {
		return nil, err
	}

	ipData := raw[14:]
	if len(ipData) < 20 {
		return nil, errIPv4HeaderTooShort
	}

	ihl := ipData[0] & 0x0F
	hLen := int(ihl) * 4
	if len(raw) < hLen {
		return nil, errIPv4HeaderLenLessThanIHL
	}

	return &IPv4Packet{
		ethFrame: *frame,
		ipHeader: ipv4HeaderfromBytes(ipData[:hLen]),
	}, nil
}

func (p *IPv4Packet) Info() string {
	pv := protocolValues[p.ipHeader.protocol]
	return fmt.Sprintf(`
IPv4 packet

===============================
%s
===============================

Version: %d
Header Length: %d bytes
DSCP: %d
ECN: %d
Total Length: %d
Identification: %d
Flags: %d
Fragment Offset: %d
TTL: %d
Transport Layer Protocol: %d (%s)
Header Checksum: %d
Source IP: %s
Destination IP: %s
Options: %v
`,
		p.ethFrame.info(),
		p.ipHeader.version, p.ipHeader.ihl*4, p.ipHeader.dscp, p.ipHeader.ecn, p.ipHeader.totalLength, p.ipHeader.identification,
		p.ipHeader.flags, p.ipHeader.fragmentOffset, p.ipHeader.ttl, p.ipHeader.protocol, pv, p.ipHeader.headerChecksum,
		p.ipHeader.sourceIP, p.ipHeader.destinationIP, p.ipHeader.options,
	)
}

func ipv4HeaderfromBytes(raw []byte) IPv4Header {
	h := IPv4Header{
		version:        raw[0] >> 4,
		ihl:            raw[0] & 0x0F,
		dscp:           raw[1] >> 2,
		ecn:            raw[1] & 0x03,
		totalLength:    binary.BigEndian.Uint16(raw[2:4]),
		identification: binary.BigEndian.Uint16(raw[4:6]),
		flags:          binary.BigEndian.Uint16(raw[6:8]) >> 13,
		fragmentOffset: binary.BigEndian.Uint16(raw[6:8]) & 0x1FFF,
		ttl:            raw[8],
		protocol:       raw[9],
		headerChecksum: binary.BigEndian.Uint16(raw[10:12]),
		sourceIP:       net.IP(raw[12:16]),
		destinationIP:  net.IP(raw[16:20]),
	}

	if len(raw) > 20 {
		h.options = raw[20:]
	}
	return h
}

func Ipv6PacketFromBytes(raw []byte) (*IPv6Packet, error) {
	frame, err := EthFrameFromBytes(raw[:14])
	if err != nil {
		return nil, err
	}

	ipData := raw[14:]
	if len(ipData) < 40 {
		return nil, errInvalidIPv6Header
	}

	return &IPv6Packet{
		ethFrame: *frame,
		ipHeader: ipv6HeaderfromBytes(ipData[:40]),
	}, nil
}

func (p *IPv6Packet) Info() string {
	pv := protocolValues[p.ipHeader.nextHeader]
	return fmt.Sprintf(`
IPv6 packet

===============================
%s
===============================

Version: %d
Header Length: 40 bytes
Traffic Class: %d
Flow Label: %d
Payload Length: %d
Transport Layer Protocol: %d (%s)
Hop Limit: %d
Source IP: %s
Destination IP: %s
`,
		p.ethFrame.info(),
		p.ipHeader.version, p.ipHeader.trafficClass, p.ipHeader.flowLabel, p.ipHeader.payloadLength,
		p.ipHeader.nextHeader, pv, p.ipHeader.hopLimit, p.ipHeader.sourceIP, p.ipHeader.destinationIP,
	)
}

func ipv6HeaderfromBytes(raw []byte) IPv6Header {
	return IPv6Header{
		version:       raw[0] >> 4,
		trafficClass:  (raw[0]&0x0F)<<4 | raw[1]>>4,
		flowLabel:     uint32(raw[1]&0x0F)<<16 | uint32(raw[2])<<8 | uint32(raw[3]),
		payloadLength: binary.BigEndian.Uint16(raw[4:6]),
		nextHeader:    raw[6],
		hopLimit:      raw[7],
		sourceIP:      net.IP(raw[8:24]),
		destinationIP: net.IP(raw[24:40]),
	}
}
