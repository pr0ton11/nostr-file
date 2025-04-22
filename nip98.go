package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/nbd-wtf/go-nostr"
)

// Event represents a Nostr event.
// Replicated from the nostr-go library because it is missing json tags.
type Event struct {
	ID        string          `json:"id"`
	PubKey    string          `json:"pubkey"`
	CreatedAt nostr.Timestamp `json:"created_at"`
	Kind      int             `json:"kind"`
	Tags      nostr.Tags      `json:"tags"`
	Content   string          `json:"content"`
	Sig       string          `json:"sig"`
}

// ValidateSignature validates the signature of the event.
// Returns true if the signature is valid, false otherwise.
// Returns an error if the public key or signature is invalid hex.
func (e Event) ValidateSignature() (bool, error) {

	pk, err := hex.DecodeString(e.PubKey)
	if err != nil {
		return false, fmt.Errorf("event pubkey '%s' is invalid hex: %w", e.PubKey, err)
	}

	pubkey, err := schnorr.ParsePubKey(pk)
	if err != nil {
		return false, fmt.Errorf("event has invalid pubkey '%s': %w", e.PubKey, err)
	}

	s, err := hex.DecodeString(e.Sig)
	if err != nil {
		return false, fmt.Errorf("signature '%s' is invalid hex: %w", e.Sig, err)
	}
	sig, err := schnorr.ParseSignature(s)
	if err != nil {
		return false, fmt.Errorf("failed to parse signature: %w", err)
	}

	hash := sha256.Sum256(e.Serialize())
	return sig.Verify(hash[:], pubkey), nil
}

// Sign signs the event with the given secret key.
// Sets the ID, PubKey and Sig fields of the event.
// Returns an error if the secret key is invalid or if signing fails.
func (e *Event) Sign(secretKey string) error {

	s, err := hex.DecodeString(secretKey)
	if err != nil {
		return fmt.Errorf("Sign called with invalid secret key '%s': %w", secretKey, err)
	}

	if e.Tags == nil {
		e.Tags = make(nostr.Tags, 0)
	}

	sk, pk := btcec.PrivKeyFromBytes(s)
	pkBytes := pk.SerializeCompressed()
	e.PubKey = hex.EncodeToString(pkBytes[1:])

	h := sha256.Sum256(e.Serialize())
	sig, err := schnorr.Sign(sk, h[:], schnorr.FastSign())
	if err != nil {
		return err
	}

	e.ID = hex.EncodeToString(h[:])
	e.Sig = hex.EncodeToString(sig.Serialize())

	return nil
}

// Serialize outputs a byte array that can be hashed to produce the canonical event "id".
func (e *Event) Serialize() []byte {
	// For some reason the content of the event is put into a byte array
	// This is done to have a predifined order of the fields
	// Specified in the NIP-01
	result := make([]byte, 0, 100+len(e.Content)+len(e.Tags)*80)
	return serializeEventInto(e, result)
}

// serializeEventInto serializes the event into a byte array.
func serializeEventInto(e *Event, result []byte) []byte {
	result = append(result, "[0,\""...)
	result = append(result, e.PubKey...)
	result = append(result, "\","...)
	result = append(result, strconv.FormatInt(int64(e.CreatedAt), 10)...)
	result = append(result, ',')
	result = append(result, strconv.Itoa(e.Kind)...)
	result = append(result, ',')

	// tags
	result = append(result, '[')
	for i, tag := range e.Tags {
		if i > 0 {
			result = append(result, ',')
		}
		// tag item
		result = append(result, '[')
		for i, s := range tag {
			if i > 0 {
				result = append(result, ',')
			}
			result = escapeString(result, s)
		}
		result = append(result, ']')
	}
	result = append(result, "],"...)

	// content needs to be escaped in general as it is user generated.
	result = escapeString(result, e.Content)
	result = append(result, ']')

	return result
}

// ValidateNIP98Event validates if the event is a NIP98 event.
// See https://github.com/nostr-protocol/nips/blob/master/98.md for specification.
func (e Event) ValidateNIP98Event(uri string, method string) bool {
	// Check if the kind of the event is 27235
	if e.Kind != 27235 {
		return false
	}
	// Check if the created at timestamp is within the last 60 seconds
	if e.CreatedAt < nostr.Now()-60 {
		return false
	}
	// Find the mandatory tags
	eventURISlice := e.Tags.Find("u")
	eventMethodSlice := e.Tags.Find("method")
	// Find the optional tags

	// This is not mandatory, but we will check if they are present
	// eventPayloadSlice := e.Tags.Find("payload")

	// For some reason the tags are modelled as a slice instead of a map
	// Key is the first element of the tag
	// Value is the second element of the tag
	// Also this could be nil for some reason
	if eventURISlice == nil || eventMethodSlice == nil {
		return false
	}
	// Do another check to make sure the slices have the right length
	// Assumed to be 2 for key and value
	if len(eventURISlice) != 2 || len(eventMethodSlice) != 2 {
		return false
	}

	// Ensure that the signature is valid
	valid, err := e.ValidateSignature()
	if err != nil || !valid {
		return false
	}

	// This event is a valid NIP98 event
	return true
}
