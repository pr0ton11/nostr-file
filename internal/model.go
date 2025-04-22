package internal

import (
	"strings"
)

// PubKey represents a public key in NPUB string format.
type PubKey string

// Invite represents an invitation to allow a user to join the file service
type Invite struct {
	Pubkey PubKey `json:"pubkey"` // Public key of the user in hex format
}

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
