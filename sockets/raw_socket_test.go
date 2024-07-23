package sockets

import (
	"reflect"
	"testing"

	"github.com/NamelessOne91/bisturi/protocols"
)

type mockIPHeader struct {
	len         int
	source      string
	destination string
	protocol    string
}

func (h *mockIPHeader) Len() int {
	return h.len
}

func (h *mockIPHeader) Source() string {
	return h.source
}

func (h *mockIPHeader) Destination() string {
	return h.destination
}

func (h *mockIPHeader) TransportLayerProtocol() string {
	return h.protocol
}

type mockIPPacket struct {
	info    string
	version uint8
	header  protocols.IPHeader
	payload []byte
}

func (p *mockIPPacket) Info() string {
	return p.info
}

func (p *mockIPPacket) Version() uint8 {
	return p.version
}

func (p *mockIPPacket) Header() protocols.IPHeader {
	return p.header
}

func (p *mockIPPacket) Payload() []byte {
	return p.payload
}

func TestHostToNetworkShort(t *testing.T) {
	tests := []struct {
		name     string
		host     uint16
		expected uint16
	}{
		{
			name:     "Little Endian to Big Endian 1",
			host:     0xff00,
			expected: 0x00ff,
		},
		{
			name:     "Little Endian to Big Endian 2",
			host:     0xf00f,
			expected: 0x0ff0,
		},
		{
			name:     "Little Endian to Big Endian 3",
			host:     0x0ff0,
			expected: 0xf00f,
		},
		{
			name:     "Little Endian to Big Endian 4",
			host:     0xfff0,
			expected: 0xf0ff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hostToNetworkShort(tt.host); got != tt.expected {
				t.Errorf("expected %016b to become %016b: got %016b", tt.host, tt.expected, got)
			}
		})
	}
}

func TestHandleLayer4Protocol(t *testing.T) {
	tests := []struct {
		name           string
		protocol       string
		packet         protocols.IPPacket
		expectedPacket NetworkPacket
		expectedErr    error
	}{
		{
			name:        "UDP filter || UDP4 packet",
			protocol:    "udp",
			expectedErr: nil,
			packet: &mockIPPacket{
				info:    "valid UDP4 packet",
				version: 4,
				header: &mockIPHeader{
					len:         20,
					source:      "192.168.0.1",
					destination: "192.168.0.2",
					protocol:    "udp",
				},
				payload: []byte{0x04, 0xd2, 0x16, 0x2e, 0x00, 0x08, 0x00, 0x00},
			},
			expectedPacket: &protocols.UDPPacket{
				IPPacket: &mockIPPacket{
					info:    "valid UDP4 packet",
					version: 4,
					header: &mockIPHeader{
						len:         20,
						source:      "192.168.0.1",
						destination: "192.168.0.2",
						protocol:    "udp",
					},
					payload: []byte{0x04, 0xd2, 0x16, 0x2e, 0x00, 0x08, 0x00, 0x00},
				},
				Header: protocols.UDPHeader{
					SourcePort:      1234,
					DestinationPort: 5678,
					Length:          8,
					Checksum:        0,
				},
			},
		},
		{
			name:        "UDP filter || UDP6 packet",
			protocol:    "udp",
			expectedErr: nil,
			packet: &mockIPPacket{
				info:    "valid UDP6 packet",
				version: 6,
				header: &mockIPHeader{
					len:         20,
					source:      "2001:db8::1",
					destination: "2001:db8::2",
					protocol:    "udp",
				},
				payload: []byte{0x04, 0xd2, 0x16, 0x2e, 0x00, 0x08, 0x00, 0x00},
			},
			expectedPacket: &protocols.UDPPacket{
				IPPacket: &mockIPPacket{
					info:    "valid UDP6 packet",
					version: 6,
					header: &mockIPHeader{
						len:         20,
						source:      "2001:db8::1",
						destination: "2001:db8::2",
						protocol:    "udp",
					},
					payload: []byte{0x04, 0xd2, 0x16, 0x2e, 0x00, 0x08, 0x00, 0x00},
				},
				Header: protocols.UDPHeader{
					SourcePort:      1234,
					DestinationPort: 5678,
					Length:          8,
					Checksum:        0,
				},
			},
		},
		{
			name:        "TCP filter || TCP4 packet",
			protocol:    "tcp",
			expectedErr: nil,
			packet: &mockIPPacket{
				info:    "valid TCP4 packet",
				version: 4,
				header: &mockIPHeader{
					len:         20,
					source:      "192.168.0.1",
					destination: "192.168.0.2",
					protocol:    "tcp",
				},
				payload: []byte{
					0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
					0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
					0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
				},
			},
			expectedPacket: &protocols.TCPPacket{
				IPPacket: &mockIPPacket{
					info:    "valid TCP4 packet",
					version: 4,
					header: &mockIPHeader{
						len:         20,
						source:      "192.168.0.1",
						destination: "192.168.0.2",
						protocol:    "tcp",
					},
					payload: []byte{
						0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
						0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
						0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
					},
				},
				Header: protocols.TCPHeader{
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
			name:        "TCP filter || TCP6 packet",
			protocol:    "tcp",
			expectedErr: nil,
			packet: &mockIPPacket{
				info:    "valid TCP6 packet",
				version: 6,
				header: &mockIPHeader{
					len:         20,
					source:      "2001:db8::1",
					destination: "2001:db8::2",
					protocol:    "tcp",
				},
				payload: []byte{
					0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
					0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
					0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
				},
			},
			expectedPacket: &protocols.TCPPacket{
				IPPacket: &mockIPPacket{
					info:    "valid TCP6 packet",
					version: 6,
					header: &mockIPHeader{
						len:         20,
						source:      "2001:db8::1",
						destination: "2001:db8::2",
						protocol:    "tcp",
					},
					payload: []byte{
						0x00, 0x50, 0x01, 0xbb, 0x1c, 0x46, 0x6f, 0x58, 0x00, 0x00, 0x00, 0x00,
						0x70, 0x02, 0x20, 0x00, 0xe0, 0x57, 0x00, 0x00,
						0x01, 0x01, 0x08, 0x0a, 0x00, 0x00, 0x00, 0x01,
					},
				},
				Header: protocols.TCPHeader{
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
			name:        "UDP filter || error",
			protocol:    "udp",
			expectedErr: protocols.ErrInvalidUDPHeader,
			packet: &mockIPPacket{
				info:    "invalid UDP4 packet",
				version: 4,
				header: &mockIPHeader{
					len:         20,
					source:      "192.168.0.1",
					destination: "192.168.0.2",
					protocol:    "udp",
				},
				payload: []byte{0x04, 0xd2},
			},
			expectedPacket: nil,
		},
		{
			name:        "TCP filter || error",
			protocol:    "tcp",
			expectedErr: protocols.ErrTCPHeaderTooShort,
			packet: &mockIPPacket{
				info:    "invalid TCP4 packet",
				version: 4,
				header: &mockIPHeader{
					len:         20,
					source:      "192.168.0.1",
					destination: "192.168.0.2",
					protocol:    "tcp",
				},
				payload: []byte{0x04, 0xd2, 0x16, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00},
			},
			expectedPacket: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataChan := make(chan NetworkPacket)
			errChan := make(chan error)

			go handleLayer4Protocol(tt.protocol, tt.packet, dataChan, errChan)

			select {
			case np := <-dataChan:
				if !reflect.DeepEqual(np, tt.expectedPacket) {
					t.Errorf("Expected packet to be: %v - got %v", tt.expectedPacket, np)
				}
			case err := <-errChan:
				if tt.expectedErr == nil || err != tt.expectedErr {
					t.Errorf("Expected error to be: %v - got %v", tt.expectedErr, err)
				}
			}
		})
	}
}
