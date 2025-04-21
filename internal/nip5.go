package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// NIP5Response represents the response from a NIP-5 request.
type NIP5Response struct {
	// The names of the users associated with the public keys.
	Names map[string]string `json:"names"`
	// The relays associated with the public keys.
	Relays map[string][]string `json:"relays"`
}

// ExtractFQUN extracts the username and domain from a fully qualified username (FQUN).
func extractFQUN(fqun string) (string, string) {
	// Split the FQUN into username and domain
	parts := strings.Split(fqun, "@")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

// RequestNIP5 sends a request to the NIP-5 endpoint and returns the response.
// The endpoint is expected to be in the format "https://<domain>/.well-known/nostr.json".
func RequestNIP5(domain string) (*NIP5Response, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", "https://"+domain+"/.well-known/nostr.json", nil)
	if err != nil {
		return nil, err
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "nostr-file")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch NIP-5 data: %s", resp.Status)
	}

	// Decode the JSON response
	var nip5Response NIP5Response
	if err := json.NewDecoder(resp.Body).Decode(&nip5Response); err != nil {
		return nil, err
	}

	return &nip5Response, nil
}

// RequestNIP5User sends a request to the NIP-5 endpoint for a specific user.
// The endpoint is expected to be in the format "https://<domain>/.well-known/nostr.json".
func RequestNIP5User(fqun string) (*NIP5Response, error) {
	// Check if the username is in the format "name@domain"
	if !strings.Contains(fqun, "@") {
		return nil, fmt.Errorf("invalid username format: %s", fqun)
	}
	// Extract the domain from the FQUN
	username, domain := extractFQUN(fqun)
	// Check if the domain is valid
	nip5Response, err := RequestNIP5(domain)
	// Check if the request was successful
	if err != nil {
		return nil, err
	}
	// Check if the username exists in the NIP-5 response
	if _, ok := nip5Response.Names[username]; !ok {
		return nil, fmt.Errorf("user %s not found in NIP-5 response", fqun)
	}
	// Return the NIP-5 response
	return nip5Response, nil
}

// LookupNIP5UserKey sends a request to the NIP-5 endpoint for a specific user and returns the public key.
func LookupNIP5UserKey(fqun string) (string, error) {
	// Extract the username and domain from the FQUN
	username, _ := extractFQUN(fqun)
	// Request the NIP-5 response for the user
	nip5Response, err := RequestNIP5User(fqun)
	if err != nil {
		return "", err
	}
	// Extract the public key from the NIP-5 response
	pubKey, ok := nip5Response.Names[username]
	if !ok {
		return "", fmt.Errorf("user %s not found in NIP-5 response", fqun)
	}
	return pubKey, nil
}

// LookupNIP5UserRelays sends a request to the NIP-5 endpoint for a specific user and returns the relays.
func LookupNIP5UserRelays(fqun string) ([]string, error) {
	// Extract the username and domain from the FQUN
	username, _ := extractFQUN(fqun)
	// Request the NIP-5 response for the user
	nip5Response, err := RequestNIP5User(fqun)
	if err != nil {
		return []string{}, err
	}
	// Extract the relays from the NIP-5 response
	relays, ok := nip5Response.Relays[username]
	if !ok {
		return []string{}, fmt.Errorf("user %s not found in NIP-5 response", username)
	}
	return relays, nil
}
