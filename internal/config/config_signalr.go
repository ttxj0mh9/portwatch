package config

import "fmt"

// SignalRConfig holds configuration for the Azure SignalR alert handler.
type SignalRConfig struct {
	Enabled   bool   `yaml:"enabled"`
	HubURL    string `yaml:"hub_url"`
	AccessKey string `yaml:"access_key"`
}

func signalRDefaults(cfg *Config) {
	// No defaults to apply — all fields must be explicitly set when enabled.
}

func validateSignalR(cfg *Config) error {
	s := cfg.Alerts.SignalR
	if !s.Enabled {
		return nil
	}
	if s.HubURL == "" {
		return fmt.Errorf("alerts.signalr.hub_url is required when signalr is enabled")
	}
	if s.AccessKey == "" {
		return fmt.Errorf("alerts.signalr.access_key is required when signalr is enabled")
	}
	return nil
}
