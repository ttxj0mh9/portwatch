package config

import "testing"

func TestZulipDefaults_Topic(t *testing.T) {
	cfg := &Config{}
	zulipDefaults(cfg)
	if cfg.Alerts.Zulip.Topic != "portwatch" {
		t.Errorf("expected default topic %q, got %q", "portwatch", cfg.Alerts.Zulip.Topic)
	}
}

func TestZulipDefaults_DoesNotOverride(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Zulip.Topic = "alerts"
	zulipDefaults(cfg)
	if cfg.Alerts.Zulip.Topic != "alerts" {
		t.Errorf("expected topic %q to be preserved, got %q", "alerts", cfg.Alerts.Zulip.Topic)
	}
}

func TestValidateZulip_Disabled(t *testing.T) {
	cfg := &Config{}
	if err := validateZulip(cfg); err != nil {
		t.Errorf("expected no error for disabled zulip, got %v", err)
	}
}

func TestValidateZulip_MissingBaseURL(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Zulip.Enabled = true
	cfg.Alerts.Zulip.Email = "bot@example.com"
	cfg.Alerts.Zulip.APIKey = "key"
	cfg.Alerts.Zulip.Stream = "general"
	if err := validateZulip(cfg); err == nil {
		t.Error("expected error for missing base_url")
	}
}

func TestValidateZulip_MissingEmail(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Zulip.Enabled = true
	cfg.Alerts.Zulip.BaseURL = "https://zulip.example.com"
	cfg.Alerts.Zulip.APIKey = "key"
	cfg.Alerts.Zulip.Stream = "general"
	if err := validateZulip(cfg); err == nil {
		t.Error("expected error for missing email")
	}
}

func TestValidateZulip_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Zulip.Enabled = true
	cfg.Alerts.Zulip.BaseURL = "https://zulip.example.com"
	cfg.Alerts.Zulip.Email = "bot@example.com"
	cfg.Alerts.Zulip.APIKey = "secret"
	cfg.Alerts.Zulip.Stream = "general"
	if err := validateZulip(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateZulip_MissingStream(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Zulip.Enabled = true
	cfg.Alerts.Zulip.BaseURL = "https://zulip.example.com"
	cfg.Alerts.Zulip.Email = "bot@example.com"
	cfg.Alerts.Zulip.APIKey = "secret"
	if err := validateZulip(cfg); err == nil {
		t.Error("expected error for missing stream")
	}
}
