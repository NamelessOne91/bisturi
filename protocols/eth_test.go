package protocols

import "testing"

func TestEthFrameFromBytes(t *testing.T) {
	tests := []struct {
		name              string
		raw               []byte
		expectedDestMAC   string
		expectedSourceMAC string
		expectedEtherType uint16
		expectedErr       error
	}{
		{
			name:              "13 bytes result in an invalid Ethernet frame",
			raw:               []byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08},
			expectedDestMAC:   "",
			expectedSourceMAC: "",
			expectedEtherType: 0,
			expectedErr:       errInvalidETHFrame,
		},
		{
			name:              "15 bytes result in a valid Ethernet frame with IPv4",
			raw:               []byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00, 0x45},
			expectedDestMAC:   "00:1a:a0:bb:cc:dd",
			expectedSourceMAC: "00:1a:b0:cc:dd:ee",
			expectedEtherType: 0x0800,
			expectedErr:       nil,
		},
		{
			name:              "Valid Ethernet frame with IPv4",
			raw:               []byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x00},
			expectedDestMAC:   "00:1a:a0:bb:cc:dd",
			expectedSourceMAC: "00:1a:b0:cc:dd:ee",
			expectedEtherType: 0x0800,
			expectedErr:       nil,
		},
		{
			name:              "Valid Ethernet frame with ARP",
			raw:               []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x08, 0x06},
			expectedDestMAC:   "ff:ff:ff:ff:ff:ff",
			expectedSourceMAC: "00:1a:b0:cc:dd:ee",
			expectedEtherType: 0x0806,
			expectedErr:       nil,
		},
		{
			name:              "Valid Ethernet frame with IPv6",
			raw:               []byte{0x00, 0x1A, 0xA0, 0xBB, 0xCC, 0xDD, 0x00, 0x1A, 0xB0, 0xCC, 0xDD, 0xEE, 0x86, 0xDD},
			expectedDestMAC:   "00:1a:a0:bb:cc:dd",
			expectedSourceMAC: "00:1a:b0:cc:dd:ee",
			expectedEtherType: 0x86DD,
			expectedErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame, err := EthFrameFromBytes(tt.raw)
			if tt.expectedErr != err {
				t.Errorf("expected error to be: %v - got: %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil {
				if tt.expectedDestMAC != frame.destinationMAC.String() || tt.expectedSourceMAC != frame.sourceMAC.String() {
					t.Errorf(
						"expected destination and source MAC to be %s and %s - got %s and %s",
						tt.expectedDestMAC, tt.expectedSourceMAC, frame.destinationMAC.String(), frame.sourceMAC.String(),
					)
				}
				if tt.expectedEtherType != frame.etherType {
					t.Errorf("expected ethernet type to be %v - got %v", tt.expectedEtherType, frame.etherType)
				}
			}
		})
	}
}
