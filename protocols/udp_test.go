package protocols

import (
	"reflect"
	"testing"
)

func TestUDPPacketFromIPPacket(t *testing.T) {
	tests := []struct {
		name              string
		IPPacket          IPPacket
		expectedUDPPacket *UDPPacket
		expectedErr       error
	}{
		{
			name: "valid IPv4 packet with UDP payload",
			IPPacket: ipv4Packet{
				payload: []byte{0x1f, 0x90, 0x23, 0xc4, 0x00, 0x10, 0x27, 0x10},
			},
			expectedUDPPacket: &UDPPacket{
				IPPacket: ipv4Packet{
					payload: []byte{0x1f, 0x90, 0x23, 0xc4, 0x00, 0x10, 0x27, 0x10},
				},
				Header: UDPHeader{
					SourcePort:      8080,
					DestinationPort: 9156,
					Length:          16,
					Checksum:        10000,
				},
			},
			expectedErr: nil,
		},
		{
			name: "IPv4 packet with too short UDP payload",
			IPPacket: ipv4Packet{
				payload: []byte{0x1f, 0x90, 0x23},
			},
			expectedUDPPacket: nil,
			expectedErr:       errInvalidUDPHeader,
		},
		{
			name: "Valid IPv6 packet with UDP payload",
			IPPacket: ipv6Packet{
				payload: []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x08, 0x00, 0x00},
			},
			expectedUDPPacket: &UDPPacket{
				IPPacket: ipv6Packet{
					payload: []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x08, 0x00, 0x00},
				},
				Header: UDPHeader{
					SourcePort:      1,
					DestinationPort: 2,
					Length:          8,
					Checksum:        0,
				},
			},
			expectedErr: nil,
		},
		{
			name: "IPv6 packet with zero Length UDP payload",
			IPPacket: ipv6Packet{
				payload: []byte{0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x9a, 0xbc},
			},
			expectedUDPPacket: &UDPPacket{
				IPPacket: ipv6Packet{
					payload: []byte{0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x9a, 0xbc},
				},
				Header: UDPHeader{
					SourcePort:      0x1234,
					DestinationPort: 0x5678,
					Length:          0,
					Checksum:        0x9abc,
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			udp, err := UDPPacketFromIPPacket(tt.IPPacket)
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
		expectedHeader *UDPHeader
		expectedErr    error
	}{
		{
			name: "Valid UDP Header",
			raw:  []byte{0x1f, 0x90, 0x23, 0xc4, 0x00, 0x10, 0x27, 0x10},
			expectedHeader: &UDPHeader{
				SourcePort:      8080,
				DestinationPort: 9156,
				Length:          16,
				Checksum:        10000,
			},
			expectedErr: nil,
		},
		{
			name:           "Too short Header",
			raw:            []byte{0x1f, 0x90, 0x23},
			expectedHeader: nil,
			expectedErr:    errInvalidUDPHeader,
		},
		{
			name: "Minimum valid Header",
			raw:  []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x08, 0x00, 0x00},
			expectedHeader: &UDPHeader{
				SourcePort:      1,
				DestinationPort: 2,
				Length:          8,
				Checksum:        0,
			},
			expectedErr: nil,
		},
		{
			name: "Zero Length Header",
			raw:  []byte{0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x9a, 0xbc},
			expectedHeader: &UDPHeader{
				SourcePort:      0x1234,
				DestinationPort: 0x5678,
				Length:          0,
				Checksum:        0x9abc,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := UDPHeaderFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedHeader) {
				t.Errorf("expected Header to be %+v - got %+v", tt.expectedHeader, h)
			}
		})
	}
}
