package config

import "fmt"

// GrafanaConfig holds configuration for the Grafana annotation handler.
type GrafanaConfig struct {
	Enabled  bool   `yaml:"enabled"`
	BaseURL  string `yaml:"base_url"`
	APIKey   string `yaml:"api_key"`
}

// grafanaDefaults applies default values to GrafanaConfig.
func grafanaDefaults(cfg *Config) {
	// No defaults needed beyond zero values; base_url and api_key must be
	// explicitly provided by the user.
}

// validateGrafana checks that required Grafana fields are present when enabled.
func validateGrafana(cfg *Config) error {
	g := cfg.Alerts.Grafana
	if !g.Enabled {
		return nil
	}
	if g.BaseURL == "" {
		return fmt.Errorf("alerts.grafana.base_url is required when grafana is enabled")
	}
	if g.APIKey == "" {
		return fmt.Errorf("alerts.grafana.api_key is required when grafana is enabled")
	}
	return nil
}
