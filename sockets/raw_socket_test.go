package sockets

import (
	"testing"
)

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
