package main

import (
	"context"
	"log/slog"
	"slices"

	"github.com/nbd-wtf/go-nostr/nip05"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/robfig/cron/v3"
)

var instance *State = nil // Singleton instance of State

// State represents the state of the application.
// Holds the configuration and other global variables in memory for easy access as a singleton.
type State struct {
	// Context is the context for the application.
	context context.Context
	// Config holds the configuration for the application.
	Config *Config
	// AllowedUsersPubkeys holds the list of allowed users' public keys.
	AllowedUsersPubkeys []string
	// AdminUsersPubkeys holds the list of admin users' public keys.
	AdminUsersPubkeys []string
	// Instance of the cron scheduler for the application.
	Cron *cron.Cron
}

// GetState returns the singleton instance of state.
func GetState() *State {
	// Check if the instance is uninitialized
	if instance == nil {
		// Initialize the instance with default values
		instance = &State{
			Cron: cron.New(),
		}
	}
	// Return the singleton instance
	return instance
}

// Updates the NIP5 records for all users.
// Periodically called by a cron job.
func (s *State) TaskUpdateNIP5() {
	if !s.Config.Security.Authorization.UseNIP5 {
		slog.Info("NIP5 is not enabled, skipping update")
		return
	}
	for _, fqun := range s.Config.Security.Authorization.AllowedUsers {
		// Update is only needed for users that are defined as Usernames
		// Hex encoded public keys are not updated as they are written in the configuration
		if IsNostrUsername(fqun) {
			// Fetch the NIP5 record for the user
			res, _, err := nip05.Fetch(s.context, fqun)
			if err != nil {
				slog.Error("Failed to fetch NIP5 record", "user", fqun, "error", err)
				continue
			}
			// Retrieve the user from the fqun
			user, domain := GetUserAndDomain(fqun)
			// Find the public key from the NIP5 record
			// User is the first part of the domain
			pubkey, ok := res.Names[user]
			if !ok {
				slog.Error("Failed to find public key in NIP5 record", "user", user, "domain", domain)
				continue
			}
			// Check if the public key is already in the list
			if slices.Contains(s.AllowedUsersPubkeys, pubkey) {
				// Add the public key to the list
				s.AllowedUsersPubkeys = append(s.AllowedUsersPubkeys, pubkey)
				slog.Info("Added new public key from NIP5 record", "user", user, "domain", domain, "pubkey", pubkey)
			} else {
				slog.Debug("Public key already exists in the list", "user", user, "domain", domain, "pubkey", pubkey)
			}
		} else {
			// Check if the user is in npub format
			if IsNostrPubkey(fqun) {
				// Convert the npub to hex
				_, data, err := nip19.Decode(fqun)
				if err != nil {
					slog.Error("Failed to decode npub", "user", fqun, "error", err)
					continue
				}
				pubkey, ok := data.(string)
				if !ok {
					slog.Error("Failed to convert npub to hex", "user", fqun)
					continue
				}
				if !slices.Contains(s.AllowedUsersPubkeys, pubkey) {
					s.AllowedUsersPubkeys = append(s.AllowedUsersPubkeys, pubkey)
					slog.Info("Added new npub public key from configuration", "user", fqun, "pubkey", pubkey)
				}
			} else if IsNostrHexPubkey(fqun) {
				// If the public key is not already in the list, add it
				if !slices.Contains(s.AllowedUsersPubkeys, fqun) {
					s.AllowedUsersPubkeys = append(s.AllowedUsersPubkeys, fqun)
					slog.Info("Added new hex public key from configuration", "user", fqun)
				}
			} else {
				slog.Error("Invalid user format", "user", fqun)
				continue
			}
		}
		for _, fqun := range s.Config.Security.Authorization.AdminUsers {
			// Update is only needed for users that are defined as Usernames
			// Hex encoded public keys are not updated as they are written in the configuration
			if IsNostrUsername(fqun) {
				// Fetch the NIP5 record for the user
				res, _, err := nip05.Fetch(s.context, fqun)
				if err != nil {
					slog.Error("Failed to fetch NIP5 record", "user", fqun, "error", err)
					continue
				}
				// Retrieve the user from the fqun
				user, domain := GetUserAndDomain(fqun)
				// Find the public key from the NIP5 record
				// User is the first part of the domain
				pubkey, ok := res.Names[user]
				if !ok {
					slog.Error("Failed to find public key in NIP5 record", "user", user, "domain", domain)
					continue
				}
				// Check if the public key is already in the list
				if slices.Contains(s.AdminUsersPubkeys, pubkey) {
					// Add the public key to the list
					s.AdminUsersPubkeys = append(s.AdminUsersPubkeys, pubkey)
					slog.Info("Added new public key from NIP5 record", "user", user, "domain", domain, "pubkey", pubkey)
				} else {
					slog.Debug("Public key already exists in the list", "user", user, "domain", domain, "pubkey", pubkey)
				}
			} else {
				// Check if the user is in npub format
				if IsNostrPubkey(fqun) {
					// Convert the npub to hex
					_, data, err := nip19.Decode(fqun)

					if err != nil {
						slog.Error("Failed to decode npub", "user", fqun, "error", err)
						continue
					}
					pubkey, ok := data.(string)
					if !ok {
						slog.Error("Failed to convert npub to hex", "user", fqun)
						continue
					}
					if !slices.Contains(s.AdminUsersPubkeys, pubkey) {
						s.AdminUsersPubkeys = append(s.AdminUsersPubkeys, pubkey)
						slog.Info("Added new npub public key from configuration", "user", fqun, "pubkey", pubkey)
					}
				} else if IsNostrHexPubkey(fqun) {
					// If the public key is not already in the list, add it
					if !slices.Contains(s.AdminUsersPubkeys, fqun) {
						s.AdminUsersPubkeys = append(s.AdminUsersPubkeys, fqun)
						slog.Info("Added new hex public key from configuration", "user", fqun)
					}
				} else {
					slog.Error("Invalid user format", "user", fqun)
					continue
				}
			}
		}
	}

}

// Starts the cron jobs for the application.
func StartCronJobs() {
	// Ensure that the state is initialized
	s := GetState()
	// Add the cron job to update NIP5 according to the configured interval
	_, err := s.Cron.AddFunc(s.Config.Security.Authorization.NIP5CronInterval, s.TaskUpdateNIP5)
	// Check for any errors while adding the cron job
	if err != nil {
		slog.Error("Failed to add cron job", "error", err)
		return
	}
	// Start the cron scheduler
	s.Cron.Start()
}
