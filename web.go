package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

func ValiddateAuthHeader(r *http.Request) (string, bool) {
	// Check if the Authorization header is present
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}
	// Check if the Authorization header is in the format "Nostr <key>"
	if strings.HasPrefix(authHeader, "Nostr") {
		// Trim the prefix and whitespaces from the key
		b64event := strings.TrimPrefix(authHeader, "Nostr")
		b64event = strings.TrimSpace(b64event)
		// Decode the base64 event
		eventByte, err := base64.StdEncoding.DecodeString(b64event)
		if err != nil {
			return "", false
		}
		// Create the empty event
		event := Event{}
		// Attempt to decode the event
		err = json.Unmarshal(eventByte, &event)
		// Invalid event format
		if err != nil {
			return "", false
		}
		// Return the public key of the event and that it is valid
		return event.PubKey, true
	}
	// Authorization header is not in the correct format
	return "", false
}

// SetCORSHeaders is a utility function that sets the CORS headers for the response.
func SetCORSHeaders(w http.ResponseWriter) {
	// Set the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("X-Powered-By", "nostr-file")
}

// HandleOptions is a middleware function that handles CORS preflight requests.
func HandleOptions(w http.ResponseWriter, r *http.Request) {
	SetCORSHeaders(w)
	w.WriteHeader(http.StatusNoContent)
}

// Handles a GET request to the server.
// Returns a file from the webserver.
func HandleGet(w http.ResponseWriter, r *http.Request) {
	// Set the CORS headers
	SetCORSHeaders(w)
	w.WriteHeader(http.StatusOK)
}

// Handles a PUT request to the server.
// Writes a file to the webserver.
func HandlePut(w http.ResponseWriter, r *http.Request) {
	// Set the CORS headers
	SetCORSHeaders(w)
	w.WriteHeader(http.StatusOK)
}
