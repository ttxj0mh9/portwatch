package config

import (
	"testing"
)

func TestRocketChatDefaults_Username(t *testing.T) {
	cfg := Default()
	if cfg.Alerts.RocketChat.Username != "portwatch" {
		t.Errorf("expected default username 'portwatch', got %q", cfg.Alerts.RocketChat.Username)
	}
}

func TestRocketChatDefaults_DoesNotOverride(t *testing.T) {
	cfg := Default()
	cfg.Alerts.RocketChat.Username = "mybot"
	rocketChatDefaults(cfg)
	if cfg.Alerts.RocketChat.Username != "mybot" {
		t.Errorf("expected username 'mybot' to be preserved, got %q", cfg.Alerts.RocketChat.Username)
	}
}

func TestValidateRocketChat_Disabled(t *testing.T) {
	cfg := Default()
	cfg.Alerts.RocketChat.Enabled = false
	if err := validateRocketChat(cfg); err != nil {
		t.Errorf("unexpected error for disabled handler: %v", err)
	}
}

func TestValidateRocketChat_MissingWebhookURL(t *testing.T) {
	cfg := Default()
	cfg.Alerts.RocketChat.Enabled = true
	cfg.Alerts.RocketChat.WebhookURL = ""
	if err := validateRocketChat(cfg); err == nil {
		t.Error("expected error for missing webhook URL")
	}
}

func TestValidateRocketChat_Valid(t *testing.T) {
	cfg := Default()
	cfg.Alerts.RocketChat.Enabled = true
	cfg.Alerts.RocketChat.WebhookURL = "http://rocket.example.com/hooks/abc123"
	if err := validateRocketChat(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
