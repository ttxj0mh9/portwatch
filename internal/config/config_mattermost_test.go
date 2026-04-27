package config

import (
	"testing"
)

func TestMattermostDefaults_Username(t *testing.T) {
	cfg := &Config{}
	mattermostDefaults(cfg)
	if cfg.Alerts.Mattermost.Username != "portwatch" {
		t.Errorf("expected default username 'portwatch', got %q", cfg.Alerts.Mattermost.Username)
	}
}

func TestMattermostDefaults_DoesNotOverride(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Mattermost.Username = "mybot"
	mattermostDefaults(cfg)
	if cfg.Alerts.Mattermost.Username != "mybot" {
		t.Errorf("expected username 'mybot', got %q", cfg.Alerts.Mattermost.Username)
	}
}

func TestValidateMattermost_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Mattermost.Enabled = false
	if err := validateMattermost(cfg); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateMattermost_MissingWebhookURL(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Mattermost.Enabled = true
	cfg.Alerts.Mattermost.WebhookURL = ""
	if err := validateMattermost(cfg); err == nil {
		t.Error("expected error for missing webhook_url")
	}
}

func TestValidateMattermost_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Mattermost.Enabled = true
	cfg.Alerts.Mattermost.WebhookURL = "https://mattermost.example.com/hooks/xyz"
	if err := validateMattermost(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateMattermost_EnabledWithURL(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Mattermost.Enabled = true
	cfg.Alerts.Mattermost.WebhookURL = "https://chat.example.org/hooks/abc123"
	cfg.Alerts.Mattermost.Username = "portwatch"
	cfg.Alerts.Mattermost.IconURL = "https://example.com/icon.png"
	if err := validateMattermost(cfg); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
