package config

import "fmt"

// PushoverConfig holds configuration for the Pushover alert handler.
type PushoverConfig struct {
	Enabled  bool   `yaml:"enabled"`
	UserKey  string `yaml:"user_key"`
	APIToken string `yaml:"api_token"`
}

// pushoverDefaults applies default values to PushoverConfig.
func pushoverDefaults(cfg *Config) {
	// No defaults needed; credentials must be explicitly provided.
}

// validatePushover checks that required Pushover fields are present when enabled.
func validatePushover(cfg *Config) error {
	p := cfg.Alerts.Pushover
	if !p.Enabled {
		return nil
	}
	if p.UserKey == "" {
		return fmt.Errorf("alerts.pushover.user_key is required when pushover is enabled")
	}
	if p.APIToken == "" {
		return fmt.Errorf("alerts.pushover.api_token is required when pushover is enabled")
	}
	return nil
}
