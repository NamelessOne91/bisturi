package protocols

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// maps the IP header values to the corresponding transport layer protocol
var protocolValues = map[uint8]string{
	1:  "icmp",
	2:  "igmp",
	6:  "tcp",
	17: "udp",
	41: "ipv6",
	58: "icmpv6",
	89: "ospf",
}

// IPPacket defines method supported both for IPv4 and IPv6 packets
type IPPacket interface {
	Info() string
	Version() uint8
	Header() IPHeader
	Payload() []byte
}

// IPHeader defines method supported both IPv4 and IPv6 headers
type IPHeader interface {
	Len() int
	Source() string
	Destination() string
	TransportLayerProtocol() string
}

// ipv4Packet contains the IP packet data (headers and payload)
type ipv4Packet struct {
	ethFrame EthernetFrame
	header   ipv4Header
	payload  []byte
}

// ipv6Packet contains the IP packet data (headers and payload)
type ipv6Packet struct {
	ethFrame EthernetFrame
	header   ipv6Header
	payload  []byte
}

type ipv4Header struct {
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

type ipv6Header struct {
	version       uint8
	trafficClass  uint8
	flowLabel     uint32
	payloadLength uint16
	nextHeader    uint8
	hopLimit      uint8
	sourceIP      net.IP
	destinationIP net.IP
}

var errInvalidIPPacket = errors.New("invalid IP packet")
var errInvalidIPVersion = errors.New("invalid IP version (not 4 or 6)")
var errIPv4HeaderTooShort = errors.New("IPv4 header must be at least 20 bytes")
var errIPv4HeaderLenLessThanIHL = errors.New("IPv4 header length less than indicated IHL")
var errInvalidIPv6Header = errors.New("IPv6 header must be 40 bytes")

// IPPacketFromBytes parses an IPv4 or IPv6 packet's data from the passed raw data and return the interface representing it.
// An error is returned if the headers' constraints are not respected.
func IPPacketFromBytes(raw []byte) (IPPacket, error) {
	if len(raw) < 15 {
		return nil, errInvalidIPPacket
	}
	version := raw[14] >> 4

	var packet IPPacket
	var err error
	if version == 4 {
		packet, err = ipv4PacketFromBytes(raw)
	} else if version == 6 {
		packet, err = ipv6PacketFromBytes(raw)
	} else {
		return nil, errInvalidIPVersion
	}

	return packet, err
}

// ipv4PacketsFromFromBytes parses an array of bytes to extract headers and payload, returning a struct pointer.
func ipv4PacketFromBytes(raw []byte) (*ipv4Packet, error) {
	frame, err := EthFrameFromBytes(raw)
	if err != nil {
		return nil, err
	}

	ipData := raw[14:]
	h, err := ipv4HeaderFromBytes(ipData)
	if err != nil {
		return nil, err
	}

	return &ipv4Packet{
		ethFrame: *frame,
		header:   *h,
		payload:  ipData[h.Len():],
	}, nil
}

// TransportLayerProtocol returns the OSI Layer 4 procotol defined in the packet's header
func (h ipv4Header) TransportLayerProtocol() string {
	return protocolValues[h.protocol]
}

// HeaderLen returns the IPv4 header length in bytes
func (h ipv4Header) Len() int {
	return int(h.ihl) * 4
}

func (p ipv4Header) Source() string {
	return p.sourceIP.String()
}

func (p ipv4Header) Destination() string {
	return p.destinationIP.String()
}

func (p ipv4Packet) Header() IPHeader {
	return p.header
}

func (p ipv4Packet) Version() uint8 {
	return p.header.version
}

func (p ipv4Packet) Payload() []byte {
	return p.payload
}

// Info returns an human-readable string containing the main IPv4 packet data
func (p ipv4Packet) Info() string {
	return fmt.Sprintf("%s IPv4 packet from IP %s to IP %s",
		p.header.TransportLayerProtocol(), p.header.sourceIP, p.header.destinationIP,
	)
}

// ipv4HeaderFromBytes parses the passed bytes to a struct containing the IP header data and returns a pointer to it.
// It expects an array of at least 20 bytes or the defined IHL
func ipv4HeaderFromBytes(raw []byte) (*ipv4Header, error) {
	if len(raw) < 20 {
		return nil, errIPv4HeaderTooShort
	}

	ihl := raw[0] & 0x0F
	hLen := int(ihl) * 4
	if len(raw) < hLen {
		return nil, errIPv4HeaderLenLessThanIHL
	}

	h := &ipv4Header{
		version:        raw[0] >> 4,
		ihl:            ihl,
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

	if hLen > 20 {
		h.options = raw[20:hLen]
	}
	return h, nil
}

// ipv6PacketsFromFromBytes parses a slice of bytes to extract headers and payload, returning a struct pointer.
func ipv6PacketFromBytes(raw []byte) (*ipv6Packet, error) {
	frame, err := EthFrameFromBytes(raw)
	if err != nil {
		return nil, err
	}

	ipData := raw[14:]
	h, err := ipv6HeaderFromBytes(ipData)
	if err != nil {
		return nil, err
	}

	return &ipv6Packet{
		ethFrame: *frame,
		header:   *h,
		payload:  ipData[h.Len():],
	}, nil
}

// TransportLayerProtocol returns the OSI Layer 4 procotol defined in the packet's header
func (h ipv6Header) TransportLayerProtocol() string {
	return protocolValues[h.nextHeader]
}

// HeaderLen returns the IPv6 header length (40) in bytes
func (p ipv6Header) Len() int {
	return 40
}

func (p ipv6Packet) Version() uint8 {
	return p.header.version
}

func (p ipv6Header) Source() string {
	return p.sourceIP.String()
}

func (p ipv6Header) Destination() string {
	return p.destinationIP.String()
}

func (p ipv6Packet) Header() IPHeader {
	return p.header
}

func (p ipv6Packet) Payload() []byte {
	return p.payload
}

// Info returns an human-readable string containing the main IPv6 packet data
func (p ipv6Packet) Info() string {
	return fmt.Sprintf("%s IPv6 packet from IP %s to IP %s",
		p.header.TransportLayerProtocol(), p.header.sourceIP, p.header.destinationIP,
	)
}

// ipv6HeaderFromBytes parses the passed bytes to a struct containing the IP header data and returns a pointer to it.
// It expects an array of at least 40 bytes
func ipv6HeaderFromBytes(raw []byte) (*ipv6Header, error) {
	if len(raw) < 40 {
		return nil, errInvalidIPv6Header
	}

	return &ipv6Header{
		version:       raw[0] >> 4,
		trafficClass:  (raw[0]&0x0F)<<4 | raw[1]>>4,
		flowLabel:     uint32(raw[1]&0x0F)<<16 | uint32(raw[2])<<8 | uint32(raw[3]),
		payloadLength: binary.BigEndian.Uint16(raw[4:6]),
		nextHeader:    raw[6],
		hopLimit:      raw[7],
		sourceIP:      net.IP(raw[8:24]),
		destinationIP: net.IP(raw[24:40]),
	}, nil
}
