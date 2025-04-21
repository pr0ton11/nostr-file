package internal

import (
	"encoding/base32"
	"encoding/hex"
	"log/slog"
	"strings"

	"github.com/btcsuite/btcutil/bech32"
)

// PubKey represents a public key in NPUB string format.
type PubKey string

func (p PubKey) String() string {
	return string(p)
}

// ToHex converts the NPUB string to a hex-encoded string.
func (p PubKey) ToHex() string {
	// Convert the bech32 string to a byte array
	hrd, data, err := bech32.Decode(p.String())
	if err != nil {
		slog.Error("Failed to decode npub bech32 string", "error", err)
		return ""
	}
	// Confirm the human-readable part is "npub"
	if hrd != "npub" {
		slog.Error("Invalid npub bech32 string", "hrd", hrd)
		return ""
	}
	// Convert the decoded byte array to a string
	hexData := string(data)
	// Decode this string from base32
	decodedData, err := base32.StdEncoding.WithPadding(base32.StdPadding).DecodeString(hexData)
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
	Pubkey PubKey `json:"pubkey"`
}
