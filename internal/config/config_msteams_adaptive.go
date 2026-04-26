package config

import "fmt"

// MSTeamsAdaptiveConfig holds configuration for the Microsoft Teams Adaptive Card handler.
type MSTeamsAdaptiveConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
}

// msteamsAdaptiveDefaults applies default values to MSTeamsAdaptiveConfig.
func msteamsAdaptiveDefaults(cfg *Config) {
	// No defaults needed beyond zero values; webhook URL must be user-supplied.
	_ = cfg
}

// validateMSTeamsAdaptive validates the MSTeamsAdaptiveConfig section.
func validateMSTeamsAdaptive(cfg *Config) error {
	if !cfg.Alerts.MSTeamsAdaptive.Enabled {
		return nil
	}
	if cfg.Alerts.MSTeamsAdaptive.WebhookURL == "" {
		return fmt.Errorf("alerts.msteams_adaptive.webhook_url is required when msteams_adaptive is enabled")
	}
	return nil
}
