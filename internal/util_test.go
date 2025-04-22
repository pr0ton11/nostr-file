package internal_test

import (
	"testing"

	"github.com/pr0ton11/nostr-file/internal"
	"github.com/stretchr/testify/assert"
)

// TestPublicKeyHex32Conversion tests the conversion of NPUB keys to hex format.
func TestPublicKeyHex32Conversion(t *testing.T) {
	tests := []struct {
		name     string
		npubkey  string
		expected string
	}{
		{
			name:     "Valid npub key A",
			npubkey:  "npub1xun57ty6dufkt58qxgk6wlcwrq47ga8vdp00wuy97c9k4dapel4qt43tq3",
			expected: "37274f2c9a6f1365d0e0322da77f0e182be474ec685ef77085f60b6ab7a1cfea",
		},
		{
			name:     "Valid npub key B",
			npubkey:  "npub198lsd3rcdw932s358lgpm47w5z2qnjr6eg0elerpqzm3muftqnxq7kp807",
			expected: "29ff06c4786b8b1542343fd01dd7cea09409c87aca1f9fe46100b71df12b04cc",
		},
		{
			name:     "Valid npub key C",
			npubkey:  "npub1jy5e5245e52mluyx5n6g9r0026g4gzkcnxws87klusserf8tm08secqr2x",
			expected: "91299a2ab4cd15bff086a4f4828def5691540ad8999d03fadfe42191a4ebdbcf",
		},
		{
			name:     "Invalid npub key",
			npubkey:  "invalid_key",
			expected: "",
		},
		{
			name:     "Empty npub key",
			npubkey:  "",
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.NPubKeyToHex(tt.npubkey)
			assert.Equal(t, tt.expected, result)
		})
	}
}
