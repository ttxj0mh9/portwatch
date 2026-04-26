package config

import (
	"testing"
)

func TestMSTeamsAdaptiveDefaults_NoOp(t *testing.T) {
	cfg := &Config{}
	msteamsAdaptiveDefaults(cfg)
	if cfg.Alerts.MSTeamsAdaptive.Enabled {
		t.Error("expected enabled to be false by default")
	}
	if cfg.Alerts.MSTeamsAdaptive.WebhookURL != "" {
		t.Error("expected webhook_url to be empty by default")
	}
}

func TestValidateMSTeamsAdaptive_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.MSTeamsAdaptive.Enabled = false
	if err := validateMSTeamsAdaptive(cfg); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateMSTeamsAdaptive_MissingWebhookURL(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.MSTeamsAdaptive.Enabled = true
	cfg.Alerts.MSTeamsAdaptive.WebhookURL = ""
	if err := validateMSTeamsAdaptive(cfg); err == nil {
		t.Error("expected error for missing webhook_url")
	}
}

func TestValidateMSTeamsAdaptive_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.MSTeamsAdaptive.Enabled = true
	cfg.Alerts.MSTeamsAdaptive.WebhookURL = "https://prod.example.com/webhook"
	if err := validateMSTeamsAdaptive(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateMSTeamsAdaptive_EnabledWithURL(t *testing.T) {
	w := writeTempConfig(t, `
alerts:
  msteams_adaptive:
    enabled: true
    webhook_url: "https://example.com/hook"
`)
	cfg, err := Load(w)
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if !cfg.Alerts.MSTeamsAdaptive.Enabled {
		t.Error("expected enabled=true")
	}
	if cfg.Alerts.MSTeamsAdaptive.WebhookURL != "https://example.com/hook" {
		t.Errorf("unexpected webhook_url: %q", cfg.Alerts.MSTeamsAdaptive.WebhookURL)
	}
}
