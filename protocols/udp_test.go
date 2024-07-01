package protocols

import (
	"reflect"
	"testing"
)

func TestUDPPacketFromIPPacket(t *testing.T) {
	tests := []struct {
		name              string
		ipPacket          IPPacket
		expectedUDPPacket *UDPPacket
		expectedErr       error
	}{
		{
			name: "valid IPv4 packet with UDP payload",
			ipPacket: ipv4Packet{
				payload: []byte{0x1f, 0x90, 0x23, 0xc4, 0x00, 0x10, 0x27, 0x10},
			},
			expectedUDPPacket: &UDPPacket{
				ipPacket: ipv4Packet{
					payload: []byte{0x1f, 0x90, 0x23, 0xc4, 0x00, 0x10, 0x27, 0x10},
				},
				header: udpHeader{
					sourcePort:      8080,
					destinationPort: 9156,
					length:          16,
					checksum:        10000,
				},
			},
			expectedErr: nil,
		},
		{
			name: "IPv4 packet with too short UDP payload",
			ipPacket: ipv4Packet{
				payload: []byte{0x1f, 0x90, 0x23},
			},
			expectedUDPPacket: nil,
			expectedErr:       errInvalidUDPHeader,
		},
		{
			name: "Valid IPv6 packet with UDP payload",
			ipPacket: ipv6Packet{
				payload: []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x08, 0x00, 0x00},
			},
			expectedUDPPacket: &UDPPacket{
				ipPacket: ipv6Packet{
					payload: []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x08, 0x00, 0x00},
				},
				header: udpHeader{
					sourcePort:      1,
					destinationPort: 2,
					length:          8,
					checksum:        0,
				},
			},
			expectedErr: nil,
		},
		{
			name: "IPv6 packet with zero length UDP payload",
			ipPacket: ipv6Packet{
				payload: []byte{0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x9a, 0xbc},
			},
			expectedUDPPacket: &UDPPacket{
				ipPacket: ipv6Packet{
					payload: []byte{0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x9a, 0xbc},
				},
				header: udpHeader{
					sourcePort:      0x1234,
					destinationPort: 0x5678,
					length:          0,
					checksum:        0x9abc,
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udp, err := UDPPacketFromIPPacket(tt.ipPacket)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(udp, tt.expectedUDPPacket) {
				t.Errorf("expected UDP packet to be %+v - got %+v", tt.expectedUDPPacket, udp)
			}
		})
	}
}

func TestUDPHeaderFromBytes(t *testing.T) {
	tests := []struct {
		name           string
		raw            []byte
		expectedHeader *udpHeader
		expectedErr    error
	}{
		{
			name: "Valid UDP header",
			raw:  []byte{0x1f, 0x90, 0x23, 0xc4, 0x00, 0x10, 0x27, 0x10},
			expectedHeader: &udpHeader{
				sourcePort:      8080,
				destinationPort: 9156,
				length:          16,
				checksum:        10000,
			},
			expectedErr: nil,
		},
		{
			name:           "Too short header",
			raw:            []byte{0x1f, 0x90, 0x23},
			expectedHeader: nil,
			expectedErr:    errInvalidUDPHeader,
		},
		{
			name: "Minimum valid header",
			raw:  []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x08, 0x00, 0x00},
			expectedHeader: &udpHeader{
				sourcePort:      1,
				destinationPort: 2,
				length:          8,
				checksum:        0,
			},
			expectedErr: nil,
		},
		{
			name: "Zero length header",
			raw:  []byte{0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x9a, 0xbc},
			expectedHeader: &udpHeader{
				sourcePort:      0x1234,
				destinationPort: 0x5678,
				length:          0,
				checksum:        0x9abc,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := udpHeaderFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedHeader) {
				t.Errorf("expected header to be %+v - got %+v", tt.expectedHeader, h)
			}
		})
	}
}
