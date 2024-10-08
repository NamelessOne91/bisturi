package protocols

import (
	"net"
	"reflect"
	"testing"
)

func TestIPPacketFromBytes(t *testing.T) {
	tests := []struct {
		name                   string
		raw                    []byte
		headerLen              int
		version                uint8
		transportLayerProtocol string
		expectedErr            error
	}{
		{
			name: "Valid IPv4 packet ",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00,
				// IPv4 Header
				0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
				0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
				0xc0, 0xa8, 0x00, 0x01,
			},
			headerLen:              20,
			version:                4,
			transportLayerProtocol: "tcp",
			expectedErr:            nil,
		},
		{
			name: "Valid IPv6 packet",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x86, 0xDD,
				// IPv6 Header
				0x60, 0x00, 0x00, 0x00, 0x00, 0x14, 0x11, 0x40,
				0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x00,
				0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x01,
			},
			headerLen:              40,
			version:                6,
			transportLayerProtocol: "udp",
			expectedErr:            nil,
		},
		{
			name: "Invalid IP packet",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00,
				// missing Header
			},
			headerLen:              0,
			version:                0,
			transportLayerProtocol: "",
			expectedErr:            errInvalidIPPacket,
		},
		{
			name: "Invalid IP packet version",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00,
				// Ipv4 Header with invalid version
				0x33, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
				0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
				0xc0, 0xa8, 0x00, 0x01,
			},
			headerLen:              0,
			version:                0,
			transportLayerProtocol: "",
			expectedErr:            errInvalidIPVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet, err := IPPacketFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil {
				if tt.headerLen != packet.Header().Len() {
					t.Errorf("expected IP header length to be %d - got %d", tt.headerLen, packet.Header().Len())
				}
				if tt.version != packet.Version() {
					t.Errorf("expected IP Packet version to be %d - got %d", tt.version, packet.Version())
				}
				if tt.transportLayerProtocol != packet.Header().TransportLayerProtocol() {
					t.Errorf("expected  IP Packet transport layer protocol to be %s - got %s", tt.transportLayerProtocol, packet.Header().TransportLayerProtocol())
				}
			}
		})
	}
}

func TestIPv4PacketFromBytes(t *testing.T) {
	tests := []struct {
		name           string
		raw            []byte
		expectedPacket *ipv4Packet
		expectedErr    error
	}{
		{
			name: "Valid IPv4 packet",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00,
				// IPv4 Header
				0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
				0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
				0xc0, 0xa8, 0x00, 0x01,
			},
			expectedPacket: &ipv4Packet{
				ethFrame: EthernetFrame{
					DestinationMAC: net.HardwareAddr([]byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD}),
					SourceMAC:      net.HardwareAddr([]byte{0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE}),
					EtherType:      0x0800,
					Payload: []byte{
						0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
						0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
						0xc0, 0xa8, 0x00, 0x01,
					},
				},
				header: ipv4Header{
					version:        4,
					ihl:            5,
					dscp:           0,
					ecn:            0,
					totalLength:    0x003c,
					identification: 0x1c46,
					flags:          2,
					fragmentOffset: 0,
					ttl:            0x40,
					protocol:       0x06,
					headerChecksum: 0xb1e6,
					sourceIP:       net.IP([]byte{192, 168, 0, 104}),
					destinationIP:  net.IP([]byte{192, 168, 0, 1}),
					options:        nil,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Valid IPv4 packet with options",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00,
				// IPv4 Header with options (IHL = 6, header length = 24 bytes)
				0x46, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
				0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
				0xc0, 0xa8, 0x00, 0x01,
				// Options (4 bytes)
				0x01, 0x02, 0x03, 0x04,
			},
			expectedPacket: &ipv4Packet{
				ethFrame: EthernetFrame{
					DestinationMAC: net.HardwareAddr([]byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD}),
					SourceMAC:      net.HardwareAddr([]byte{0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE}),
					EtherType:      0x0800,
					Payload: []byte{
						0x46, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
						0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
						0xc0, 0xa8, 0x00, 0x01, 0x01, 0x02, 0x03, 0x04,
					},
				},
				header: ipv4Header{
					version:        4,
					ihl:            6,
					dscp:           0,
					ecn:            0,
					totalLength:    0x003c,
					identification: 0x1c46,
					flags:          2,
					fragmentOffset: 0,
					ttl:            0x40,
					protocol:       0x06,
					headerChecksum: 0xb1e6,
					sourceIP:       net.IP([]byte{192, 168, 0, 104}),
					destinationIP:  net.IP([]byte{192, 168, 0, 1}),
					options:        []byte{0x01, 0x02, 0x03, 0x04},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Ethernet frame",
			raw: []byte{
				// Incomplete Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00,
			},
			expectedPacket: nil,
			expectedErr:    errInvalidETHFrame,
		},
		{
			name: "Invalid IPv4 header",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00,
				// Incomplete IPv4 Header
				0x45, 0x00, 0x00, 0x3c, 0x1c,
			},
			expectedPacket: nil,
			expectedErr:    errIPv4HeaderTooShort,
		},
		{
			name: "IHL does not match actual header length",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00,
				// IPv4 Header with IHL indicating 6 words (24 bytes) but actual length is only 20 bytes
				0x46, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
				0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
				0xc0, 0xa8, 0x00, 0x01,
			},
			expectedPacket: nil,
			expectedErr:    errIPv4HeaderLenLessThanIHL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := ipv4PacketFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !compareIPv4Packets(tt.expectedPacket, p) {
				t.Errorf("expected IPv4 packet to be %+v - got %+v", tt.expectedPacket, p)
			}
		})
	}
}

func TestIPv4HeaderFromBytes(t *testing.T) {
	tests := []struct {
		name           string
		raw            []byte
		expectedHeader *ipv4Header
		expectedErr    error
	}{
		{
			name:           "Empty input",
			raw:            []byte{},
			expectedHeader: nil,
			expectedErr:    errIPv4HeaderTooShort,
		},
		{
			name:           "Incomplete header (less than 20 bytes)",
			raw:            []byte{0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00, 0x40, 0x06},
			expectedHeader: nil,
			expectedErr:    errIPv4HeaderTooShort,
		},
		{
			name: "Header length less than indicated IHL",
			raw: []byte{
				0x46, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
				0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
				0xc0, 0xa8, 0x00, 0x01,
			},
			expectedHeader: nil,
			expectedErr:    errIPv4HeaderLenLessThanIHL,
		},
		{
			name: "Valid header without options",
			raw: []byte{
				0x45, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
				0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
				0xc0, 0xa8, 0x00, 0x01,
			},
			expectedHeader: &ipv4Header{
				version:        4,
				ihl:            5,
				dscp:           0,
				ecn:            0,
				totalLength:    0x003c,
				identification: 0x1c46,
				flags:          2,
				fragmentOffset: 0,
				ttl:            0x40,
				protocol:       0x06,
				headerChecksum: 0xb1e6,
				sourceIP:       net.IP([]byte{192, 168, 0, 104}),
				destinationIP:  net.IP([]byte{192, 168, 0, 1}),
				options:        nil,
			},
			expectedErr: nil,
		},
		{
			name: "Valid header with options",
			raw: []byte{
				0x46, 0x00, 0x00, 0x3c, 0x1c, 0x46, 0x40, 0x00,
				0x40, 0x06, 0xb1, 0xe6, 0xc0, 0xa8, 0x00, 0x68,
				0xc0, 0xa8, 0x00, 0x01,
				// Options (4 bytes of options)
				0x01, 0x02, 0x03, 0x04,
			},
			expectedHeader: &ipv4Header{
				version:        4,
				ihl:            6,
				dscp:           0,
				ecn:            0,
				totalLength:    0x003c,
				identification: 0x1c46,
				flags:          2,
				fragmentOffset: 0,
				ttl:            0x40,
				protocol:       0x06,
				headerChecksum: 0xb1e6,
				sourceIP:       net.IP([]byte{192, 168, 0, 104}),
				destinationIP:  net.IP([]byte{192, 168, 0, 1}),
				options:        []byte{0x01, 0x02, 0x03, 0x04},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := ipv4HeaderFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedHeader) {
				t.Errorf("expected header to be %+v - got %+v", tt.expectedHeader, h)
			}
		})
	}
}

func TestIPv6PacketFromBytes(t *testing.T) {
	tests := []struct {
		name           string
		raw            []byte
		expectedPacket *ipv6Packet
		expectedErr    error
	}{
		{
			name: "Valid IPv6 packet",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x86, 0xDD,
				// IPv6 Header
				0x60, 0x00, 0x00, 0x00, 0x00, 0x14, 0x11, 0x40,
				0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x00,
				0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x01,
			},
			expectedPacket: &ipv6Packet{
				ethFrame: EthernetFrame{
					DestinationMAC: net.HardwareAddr([]byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD}),
					SourceMAC:      net.HardwareAddr([]byte{0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE}),
					EtherType:      0x86DD,
					Payload: []byte{
						0x60, 0x00, 0x00, 0x00, 0x00, 0x14, 0x11, 0x40,
						0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x00,
						0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x01,
					},
				},
				header: ipv6Header{
					version:       6,
					trafficClass:  0,
					flowLabel:     0,
					payloadLength: 0x0014,
					nextHeader:    0x11,
					hopLimit:      0x40,
					sourceIP:      net.IP([]byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x00}),
					destinationIP: net.IP([]byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x01}),
				},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Ethernet frame",
			raw: []byte{
				// Incomplete Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00,
			},
			expectedPacket: nil,
			expectedErr:    errInvalidETHFrame,
		},
		{
			name: "Invalid IPv6 header",
			raw: []byte{
				// Ethernet Frame
				0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x86, 0xDD,
				// Incomplete IPv6 Header
				0x60, 0x00, 0x00, 0x00, 0x00, 0x14, 0x11,
			},
			expectedPacket: nil,
			expectedErr:    errInvalidIPv6Header,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := ipv6PacketFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !compareIPv6Packets(tt.expectedPacket, p) {
				t.Errorf("expected IPv6 packet to be %+v - got %+v", tt.expectedPacket, p)
			}
		})
	}
}

func TestIPv6HeaderFromBytes(t *testing.T) {
	tests := []struct {
		name           string
		raw            []byte
		expectedHeader *ipv6Header
		expectedErr    error
	}{
		{
			name:           "Empty input",
			raw:            []byte{},
			expectedHeader: nil,
			expectedErr:    errInvalidIPv6Header,
		},
		{
			name:           "Incomplete header (less than 40 bytes)",
			raw:            []byte{0x60, 0x00, 0x00, 0x00, 0x00, 0x14, 0x11},
			expectedHeader: nil,
			expectedErr:    errInvalidIPv6Header,
		},
		{
			name: "Valid header",
			raw: []byte{
				0x60, 0x00, 0x00, 0x00, 0x00, 0x14, 0x11, 0x40,
				0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x00,
				0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x01,
			},
			expectedHeader: &ipv6Header{
				version:       6,
				trafficClass:  0,
				flowLabel:     0,
				payloadLength: 0x0014,
				nextHeader:    0x11,
				hopLimit:      0x40,
				sourceIP:      net.IP([]byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x00}),
				destinationIP: net.IP([]byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x1c, 0x7e, 0xff, 0xfe, 0xe4, 0x2c, 0x01}),
			},
			expectedErr: nil,
		},
		{
			name: "Another valid header",
			raw: []byte{
				0x60, 0x00, 0x00, 0x00, 0x00, 0x14, 0x06, 0x40,
				0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00,
				0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34,
				0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00,
				0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34,
			},
			expectedHeader: &ipv6Header{
				version:       6,
				trafficClass:  0,
				flowLabel:     0,
				payloadLength: 0x0014,
				nextHeader:    0x06,
				hopLimit:      0x40,
				sourceIP:      net.IP([]byte{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}),
				destinationIP: net.IP([]byte{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}),
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := ipv6HeaderFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedHeader) {
				t.Errorf("expected header to be %+v - got %+v", tt.expectedHeader, h)
			}
		})
	}
}

func compareIPv4Packets(a, b *ipv4Packet) bool {
	return reflect.DeepEqual(a.ethFrame, b.ethFrame) &&
		a.header.version == b.header.version &&
		a.header.ihl == b.header.ihl &&
		a.header.dscp == b.header.dscp &&
		a.header.ecn == b.header.ecn &&
		a.header.totalLength == b.header.totalLength &&
		a.header.flags == b.header.flags &&
		a.header.identification == b.header.identification &&
		a.header.fragmentOffset == b.header.fragmentOffset &&
		a.header.protocol == b.header.protocol &&
		a.header.headerChecksum == b.header.headerChecksum &&
		a.header.sourceIP.Equal(b.header.sourceIP) &&
		a.header.destinationIP.Equal(b.header.destinationIP) &&
		reflect.DeepEqual(a.header.options, b.header.options)
}

func compareIPv6Packets(a, b *ipv6Packet) bool {
	return reflect.DeepEqual(a.ethFrame, b.ethFrame) &&
		a.header.version == b.header.version &&
		a.header.trafficClass == b.header.trafficClass &&
		a.header.flowLabel == b.header.flowLabel &&
		a.header.payloadLength == b.header.payloadLength &&
		a.header.nextHeader == b.header.nextHeader &&
		a.header.hopLimit == b.header.hopLimit &&
		a.header.sourceIP.Equal(b.header.sourceIP) &&
		a.header.destinationIP.Equal(b.header.destinationIP)
}
