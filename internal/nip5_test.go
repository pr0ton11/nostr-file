package internal

import "testing"

// TestNIP5Request tests the RequestNIP5 and RequestNIP5User functions.
// It checks if the functions return the expected results for different domains and usernames.
// The test cases include both valid and invalid domains and usernames.
// It also verifies if the names returned in the response match the expected names.
func TestNIP5Request(t *testing.T) {
	tests := []struct {
		name        string
		domain      string
		pubKey      string
		errExpected bool
		exists      bool
	}{
		{
			name:        "ms",
			domain:      "pr0.guru",
			pubKey:      "37274f2c9a6f1365d0e0322da77f0e182be474ec685ef77085f60b6ab7a1cfea",
			errExpected: false,
			exists:      true,
		},
		{
			name:        "ws",
			domain:      "pr0.guru",
			pubKey:      "29ff06c4786b8b1542343fd01dd7cea09409c87aca1f9fe46100b71df12b04cc",
			errExpected: false,
			exists:      true,
		},
		{
			name:        "opc",
			domain:      "opc6.net",
			pubKey:      "91299a2ab4cd15bff086a4f4828def5691540ad8999d03fadfe42191a4ebdbcf",
			errExpected: false,
			exists:      true,
		},
		{
			name:        "hello",
			domain:      "invalid.domain",
			errExpected: true,
			exists:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nip5Response, err := RequestNIP5(tt.domain)
			if (err != nil) != tt.errExpected {
				t.Errorf("RequestNIP5() error = %v, expected %v", err, tt.errExpected)
				return
			}
			if tt.exists {
				name, ok := nip5Response.Names[tt.name]
				if !ok {
					t.Errorf("RequestNIP5() name = %v, expected %v", name, tt.name)
				}
			}
			// Build the full username to check the RequestNIP5User function
			username := tt.name + "@" + tt.domain
			nip5ResponseUser, err := RequestNIP5User(username)
			if (err != nil) != tt.errExpected {
				t.Errorf("RequestNIP5User() error = %v, expected %v", err, tt.errExpected)
				return
			}
			if tt.exists {
				name, ok := nip5ResponseUser.Names[tt.name]
				if !ok {
					t.Errorf("RequestNIP5User() name = %v, expected %v", name, tt.name)
				}
			}
			// Check the LookupNIP5UserKey function
			pubKey, err := LookupNIP5UserKey(username)
			if (err != nil) != tt.errExpected {
				t.Errorf("LookupNIP5UserKey() error = %v, expected %v", err, tt.errExpected)
				return
			}
			if tt.exists {
				if pubKey != tt.pubKey {
					t.Errorf("LookupNIP5UserKey() pubKey = %v, expected %v", pubKey, tt.pubKey)
				}
			}
		})
	}
}
