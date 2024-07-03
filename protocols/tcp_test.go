package protocols

import (
	"reflect"
	"testing"
)

func TestTCPPacketFromIPPacket(t *testing.T) {
	tests := []struct {
		name              string
		ipPacket          IPPacket
		expectedErr       error
		expectedTCPPacket *TCPPacket
	}{
		{
			name: "IPv4 Packet - Valid",
			ipPacket: ipv4Packet{
				payload: []byte{
					0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
					0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
					0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
				},
			},
			expectedErr: nil,
			expectedTCPPacket: &TCPPacket{
				ipPacket: ipv4Packet{
					payload: []byte{
						0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
						0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
						0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
					},
				},
				header: tcpHeader{
					sourcePort:      80,
					destinationPort: 443,
					sequenceNumber:  474378072,
					ackNumber:       0,
					rawOffset:       7,
					flags:           0x2,
					windowSize:      8192,
					checksum:        0xe057,
					urgentPointer:   0,
					options:         []byte{0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01},
				},
			},
		},
		{
			name: "IPv6 Packet - Valid",
			ipPacket: ipv6Packet{
				payload: []byte{
					0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
					0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
					0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
				},
			},
			expectedErr: nil,
			expectedTCPPacket: &TCPPacket{
				ipPacket: ipv6Packet{
					payload: []byte{
						0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
						0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
						0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
					},
				},
				header: tcpHeader{
					sourcePort:      80,
					destinationPort: 443,
					sequenceNumber:  474378072,
					ackNumber:       0,
					rawOffset:       7,
					flags:           0x2,
					windowSize:      8192,
					checksum:        0xe057,
					urgentPointer:   0,
					options:         []byte{0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcp, err := TCPPacketFromIPPacket(tt.ipPacket)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(tcp, tt.expectedTCPPacket) {
				t.Errorf("expected TCP packet to be %+v - got %+v", tt.expectedTCPPacket, tcp)
			}
		})
	}
}

func TestTCPHeaderFromBytes(t *testing.T) {
	tests := []struct {
		name           string
		raw            []byte
		expectedHeader *tcpHeader
		expectedErr    error
	}{
		{
			name: "Valid TCP header without options",
			raw: []byte{
				0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
				0x50, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
			},
			expectedHeader: &tcpHeader{
				sourcePort:      80,
				destinationPort: 443,
				sequenceNumber:  474378072,
				ackNumber:       0,
				rawOffset:       5,
				flags:           0x2,
				windowSize:      8192,
				checksum:        0xe057,
				urgentPointer:   0,
				options:         []byte{},
			},
			expectedErr: nil,
		},
		{
			name: "Valid TCP header with options",
			raw: []byte{
				0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
				0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
				0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
			},
			expectedHeader: &tcpHeader{
				sourcePort:      80,
				destinationPort: 443,
				sequenceNumber:  474378072,
				ackNumber:       0,
				rawOffset:       7,
				flags:           0x2,
				windowSize:      8192,
				checksum:        0xe057,
				urgentPointer:   0,
				options:         []byte{0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01},
			},
			expectedErr: nil,
		},
		{
			name:           "Invalid TCP header (too short)",
			raw:            []byte{0x00, 0x50},
			expectedHeader: nil,
			expectedErr:    errTCPHeaderTooShort,
		},
		{
			name: "Valid TCP header with minimum length",
			raw: []byte{
				0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
				0x50, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
			},
			expectedHeader: &tcpHeader{
				sourcePort:      80,
				destinationPort: 443,
				sequenceNumber:  474378072,
				ackNumber:       0,
				rawOffset:       5,
				flags:           0x2,
				windowSize:      8192,
				checksum:        0xe057,
				urgentPointer:   0,
				options:         []byte{},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid TCP header (length mismatch)",
			raw: []byte{
				0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
				0x80, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
			},
			expectedHeader: nil,
			expectedErr:    errTCPHeaderLenMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := tcpHeaderFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedHeader) {
				t.Errorf("expected header to be %+v - got %+v", tt.expectedHeader, h)
			}
		})
	}
}
