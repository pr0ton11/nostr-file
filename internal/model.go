package internal

import (
	"encoding/base32"
	"encoding/hex"
	"log/slog"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
)

// B32AllowedChars is the set of allowed characters for base32 encoding.
var B32AllowedChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// PubKey represents a public key in NPUB string format.
type PubKey string

// String() converts the PubKey to a string.
func (p PubKey) String() string {
	return string(p)
}

// Returns a new PubKey from a hex string.
// Automatically converts the npub string to hex format.
func NewPubKey(key string) PubKey {
	if strings.HasPrefix(key, "npub") {
		key = NPubKeyToHex(key)
	}
	return PubKey(key)
}

// ToHex converts the NPUB string to a hex-encoded string.
func NPubKeyToHex(NPubKey string) string {
	// Convert the bech32 string to a byte array
	hrd, data, err := bech32.Decode(NPubKey)
	if err != nil {
		slog.Error("Failed to decode npub bech32 string", "error", err)
		return ""
	}
	// Confirm the human-readable part is "npub"
	if hrd != "npub" {
		slog.Error("Invalid npub bech32 string", "hrd", hrd)
		return ""
	}
	var allowedCharsDataString string
	// Build the string from allowed character indexes
	for i := range data {
		// Convert the byte to a character
		allowedCharsDataString += string(B32AllowedChars[data[i]])
	}
	// Decode this string from base32
	decodedData, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(allowedCharsDataString)
	if err != nil {
		slog.Error("Failed to decode base32 string with padding", "error", err)
		return ""
	}
	// Convert the decoded byte array to a hex string
	b16encodedData := hex.EncodeToString(decodedData)
	// Return the hex-encoded string
	return strings.ToLower(b16encodedData)
}

// Invite represents an invitation to allow a user to join the file service
type Invite struct {
	Pubkey PubKey `json:"pubkey"` // Public key of the user in hex format
}
