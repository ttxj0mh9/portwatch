package config

import "fmt"

// NewRelicConfig holds configuration for the New Relic Insights alert handler.
type NewRelicConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AccountID string `yaml:"account_id"`
	APIKey    string `yaml:"api_key"`
}

func newRelicDefaults(cfg *Config) {
	// No default values required; credentials must be explicitly provided.
}

func validateNewRelic(cfg *Config) error {
	nr := cfg.Alerts.NewRelic
	if !nr.Enabled {
		return nil
	}
	if nr.AccountID == "" {
		return fmt.Errorf("alerts.newrelic.account_id is required when New Relic is enabled")
	}
	if nr.APIKey == "" {
		return fmt.Errorf("alerts.newrelic.api_key is required when New Relic is enabled")
	}
	return nil
}
