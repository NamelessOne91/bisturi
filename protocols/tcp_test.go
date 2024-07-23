package protocols

import (
	"reflect"
	"testing"
)

func TestTCPPacketFromIPPacket(t *testing.T) {
	tests := []struct {
		name              string
		IPPacket          IPPacket
		expectedErr       error
		expectedTCPPacket *TCPPacket
	}{
		{
			name: "IPv4 Packet - Valid",
			IPPacket: ipv4Packet{
				payload: []byte{
					0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
					0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
					0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
				},
			},
			expectedErr: nil,
			expectedTCPPacket: &TCPPacket{
				IPPacket: ipv4Packet{
					payload: []byte{
						0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
						0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
						0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
					},
				},
				Header: TCPHeader{
					SourcePort:      80,
					DestinationPort: 443,
					SequenceNumber:  474378072,
					AckNumber:       0,
					RawOffset:       7,
					Flags:           0x2,
					WindowSize:      8192,
					Checksum:        0xe057,
					UrgentPointer:   0,
					Options:         []byte{0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01},
				},
			},
		},
		{
			name: "IPv6 Packet - Valid",
			IPPacket: ipv6Packet{
				payload: []byte{
					0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
					0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
					0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
				},
			},
			expectedErr: nil,
			expectedTCPPacket: &TCPPacket{
				IPPacket: ipv6Packet{
					payload: []byte{
						0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
						0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
						0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
					},
				},
				Header: TCPHeader{
					SourcePort:      80,
					DestinationPort: 443,
					SequenceNumber:  474378072,
					AckNumber:       0,
					RawOffset:       7,
					Flags:           0x2,
					WindowSize:      8192,
					Checksum:        0xe057,
					UrgentPointer:   0,
					Options:         []byte{0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcp, err := TCPPacketFromIPPacket(tt.IPPacket)
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
		expectedHeader *TCPHeader
		expectedErr    error
	}{
		{
			name: "Valid TCP Header without Options",
			raw: []byte{
				0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
				0x50, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
			},
			expectedHeader: &TCPHeader{
				SourcePort:      80,
				DestinationPort: 443,
				SequenceNumber:  474378072,
				AckNumber:       0,
				RawOffset:       5,
				Flags:           0x2,
				WindowSize:      8192,
				Checksum:        0xe057,
				UrgentPointer:   0,
				Options:         nil,
			},
			expectedErr: nil,
		},
		{
			name: "Valid TCP Header with Options",
			raw: []byte{
				0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
				0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
				0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
			},
			expectedHeader: &TCPHeader{
				SourcePort:      80,
				DestinationPort: 443,
				SequenceNumber:  474378072,
				AckNumber:       0,
				RawOffset:       7,
				Flags:           0x2,
				WindowSize:      8192,
				Checksum:        0xe057,
				UrgentPointer:   0,
				Options:         []byte{0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01},
			},
			expectedErr: nil,
		},
		{
			name:           "Invalid TCP Header (too short)",
			raw:            []byte{0x00, 0x50},
			expectedHeader: nil,
			expectedErr:    ErrTCPHeaderTooShort,
		},
		{
			name: "Valid TCP Header with minimum length",
			raw: []byte{
				0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
				0x50, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
			},
			expectedHeader: &TCPHeader{
				SourcePort:      80,
				DestinationPort: 443,
				SequenceNumber:  474378072,
				AckNumber:       0,
				RawOffset:       5,
				Flags:           0x2,
				WindowSize:      8192,
				Checksum:        0xe057,
				UrgentPointer:   0,
				Options:         nil,
			},
			expectedErr: nil,
		},
		{
			name: "Invalid TCP Header (length mismatch)",
			raw: []byte{
				0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
				0x80, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
			},
			expectedHeader: nil,
			expectedErr:    ErrTCPHeaderLenMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := TCPHeaderFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error: %v - got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(h, tt.expectedHeader) {
				t.Errorf("expected Header to be %+v - got %+v", tt.expectedHeader, h)
			}
		})
	}
}
