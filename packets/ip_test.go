package packets

import (
	"net"
	"reflect"
	"testing"
)

func TestIPv4PacketFromBytes(t *testing.T) {
	tests := []struct {
		name           string
		raw            []byte
		expectedPacket *IPv4Packet
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
			expectedPacket: &IPv4Packet{
				ethFrame: EthernetFrame{
					destinationMAC: net.HardwareAddr([]byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD}),
					sourceMAC:      net.HardwareAddr([]byte{0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE}),
					etherType:      0x0800,
				},
				ipHeader: ipv4Header{
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
			expectedPacket: &IPv4Packet{
				ethFrame: EthernetFrame{
					destinationMAC: net.HardwareAddr([]byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD}),
					sourceMAC:      net.HardwareAddr([]byte{0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE}),
					etherType:      0x0800,
				},
				ipHeader: ipv4Header{
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
			h, err := IPv4PacketFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedPacket) {
				t.Errorf("expected IPv4 packet to be %+v - got %+v", tt.expectedPacket, h)
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
			h, err := ipv4HeaderfromBytes(tt.raw)
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
		expectedPacket *IPv6Packet
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
			expectedPacket: &IPv6Packet{
				ethFrame: EthernetFrame{
					destinationMAC: net.HardwareAddr([]byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD}),
					sourceMAC:      net.HardwareAddr([]byte{0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE}),
					etherType:      0x86DD,
				},
				ipHeader: ipv6Header{
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
			h, err := IPv6PacketFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedPacket) {
				t.Errorf("expected IPv6 packet to be %+v - got %+v", tt.expectedPacket, h)
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
			h, err := ipv6HeaderfromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedHeader) {
				t.Errorf("expected header to be %+v - got %+v", tt.expectedHeader, h)
			}
		})
	}
}
