package config

import "fmt"

// NtfyConfig holds configuration for the ntfy alert handler.
type NtfyConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ServerURL string `yaml:"server_url"`
	Topic     string `yaml:"topic"`
}

// ntfyDefaults applies default values to NtfyConfig fields.
func ntfyDefaults(cfg *Config) {
	if cfg.Alerts.Ntfy.ServerURL == "" {
		cfg.Alerts.Ntfy.ServerURL = "https://ntfy.sh"
	}
}

// validateNtfy checks that the NtfyConfig is valid when enabled.
func validateNtfy(cfg *Config) error {
	n := cfg.Alerts.Ntfy
	if !n.Enabled {
		return nil
	}
	if n.ServerURL == "" {
		return fmt.Errorf("alerts.ntfy.server_url is required when ntfy is enabled")
	}
	if n.Topic == "" {
		return fmt.Errorf("alerts.ntfy.topic is required when ntfy is enabled")
	}
	return nil
}
