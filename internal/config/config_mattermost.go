package config

import "fmt"

// MattermostConfig holds settings for the Mattermost alert handler.
type MattermostConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
	Username   string `yaml:"username"`
	IconURL    string `yaml:"icon_url"`
}

// mattermostDefaults applies default values to MattermostConfig fields.
func mattermostDefaults(cfg *Config) {
	if cfg.Alerts.Mattermost.Username == "" {
		cfg.Alerts.Mattermost.Username = "portwatch"
	}
}

// validateMattermost checks that required Mattermost fields are present when enabled.
func validateMattermost(cfg *Config) error {
	m := cfg.Alerts.Mattermost
	if !m.Enabled {
		return nil
	}
	if m.WebhookURL == "" {
		return fmt.Errorf("alerts.mattermost.webhook_url is required when mattermost is enabled")
	}
	return nil
}
