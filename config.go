package main

type Config struct {
	// The port to listen on
	Port string `yaml:"port"`
	// The path to the directory to serve files from
	StoragePath string `yaml:"storagePath"`

	Security struct {
		Authorization struct {
			Enabled          bool     `yaml:"enabled"`
			Header           string   `yaml:"header"`
			AllowedUsers     []string `yaml:"allowedUsers"`
			AdminUsers       []string `yaml:"adminUsers"`
			UseNIP5          bool     `yaml:"useNIP5"`
			NIP5CronInterval string   `yaml:"nip5CronInterval"`
		} `yaml:"authorization"`
		EnableCORS     bool     `yaml:"enableCORS"`
		AllowedOrigins []string `yaml:"allowedOrigins"`
	} `yaml:"security"`
}
