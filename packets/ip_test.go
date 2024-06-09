package packets

import (
	"net"
	"reflect"
	"testing"
)

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
