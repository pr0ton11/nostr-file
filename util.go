package main

import (
	"strings"

	"github.com/nbd-wtf/go-nostr/nip19"
)

// IsNostrUsername checks if a given string is a valid Nostr username.
// Format is: <username>@<domain>
func IsNostrUsername(username string) bool {
	// Check if the username contains an '@' symbol
	if !strings.Contains(username, "@") {
		return false
	}

	// Split the username into parts
	parts := strings.Split(username, "@")
	if len(parts) != 2 {
		return false
	}

	// Check if the domain part contains a dot
	// As every domain must have a tld
	if !strings.Contains(parts[1], ".") {
		return false
	}

	return true
}

// GetUserAndDomain extracts the username and domain from a Nostr username.
func GetUserAndDomain(username string) (string, string) {
	// Validate the username format
	if !IsNostrUsername(username) {
		return "", ""
	}
	// Split the username into parts
	parts := strings.Split(username, "@")
	// Return the username and domain
	// The username is the first part and the domain is the second part
	return parts[0], parts[1]
}

// IsNostrPubkey checks if a given string is a valid Nostr public key.
func IsNostrPubkey(key string) bool {
	// Decode the key using nip19
	// This will check if the key is a valid Nostr public key
	prefix, _, err := nip19.Decode(key)
	if err != nil {
		return false
	}
	// Check if the prefix is "npub"
	if prefix != "npub" {
		return false
	}
	return true
}

// IsNostrHexPubkey checks if a given string is a valid Nostr public key in hex format.
func IsNostrHexPubkey(key string) bool {
	// Attempt to encode the public key using nip19
	// This will check if the key is a valid Nostr public key
	_, err := nip19.EncodePublicKey(key)
	return err == nil
}

// Escaping strings for JSON encoding according to RFC8259.
// Also encloses result in quotation marks "".
func escapeString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			dst = append(dst, []byte{'\\', '"'}...)
		case c == '\\':
			dst = append(dst, []byte{'\\', '\\'}...)
		case c >= 0x20:
			dst = append(dst, c)
		case c == 0x08:
			dst = append(dst, []byte{'\\', 'b'}...)
		case c < 0x09:
			dst = append(dst, []byte{'\\', 'u', '0', '0', '0', '0' + c}...)
		case c == 0x09:
			dst = append(dst, []byte{'\\', 't'}...)
		case c == 0x0a:
			dst = append(dst, []byte{'\\', 'n'}...)
		case c == 0x0c:
			dst = append(dst, []byte{'\\', 'f'}...)
		case c == 0x0d:
			dst = append(dst, []byte{'\\', 'r'}...)
		case c < 0x10:
			dst = append(dst, []byte{'\\', 'u', '0', '0', '0', 0x57 + c}...)
		case c < 0x1a:
			dst = append(dst, []byte{'\\', 'u', '0', '0', '1', 0x20 + c}...)
		case c < 0x20:
			dst = append(dst, []byte{'\\', 'u', '0', '0', '1', 0x47 + c}...)
		}
	}
	dst = append(dst, '"')
	return dst
}
