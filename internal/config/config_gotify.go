package config

import "fmt"

type GotifyConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ServerURL string `yaml:"server_url"`
	Token     string `yaml:"token"`
}

func gotifyDefaults(cfg *Config) {
	// No computed defaults required; fields are zero-valued when absent.
}

func validateGotify(cfg *Config) error {
	g := cfg.Gotify
	if !g.Enabled {
		return nil
	}
	if g.ServerURL == "" {
		return fmt.Errorf("gotify: server_url is required when enabled")
	}
	if g.Token == "" {
		return fmt.Errorf("gotify: token is required when enabled")
	}
	return nil
}
