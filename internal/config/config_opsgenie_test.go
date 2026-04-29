package config

import (
	"testing"
)

func TestOpsGenieDefaults_APIURL(t *testing.T) {
	cfg := &Config{}
	opsgenieDefaults(cfg)
	if cfg.Alerts.OpsGenie.APIURL != "https://api.opsgenie.com/v2/alerts" {
		t.Errorf("expected default api_url, got %q", cfg.Alerts.OpsGenie.APIURL)
	}
}

func TestOpsGenieDefaults_DoesNotOverride(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.OpsGenie.APIURL = "https://custom.example.com/alerts"
	opsgenieDefaults(cfg)
	if cfg.Alerts.OpsGenie.APIURL != "https://custom.example.com/alerts" {
		t.Errorf("expected custom api_url to be preserved, got %q", cfg.Alerts.OpsGenie.APIURL)
	}
}

func TestValidateOpsGenie_Disabled(t *testing.T) {
	cfg := &Config{}
	if err := validateOpsGenie(cfg); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateOpsGenie_MissingAPIKey(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.OpsGenie.Enabled = true
	if err := validateOpsGenie(cfg); err == nil {
		t.Error("expected error for missing api_key")
	}
}

func TestValidateOpsGenie_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.OpsGenie.Enabled = true
	cfg.Alerts.OpsGenie.APIKey = "test-key-123"
	if err := validateOpsGenie(cfg); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}

func TestValidateOpsGenie_EnabledWithTags(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.OpsGenie.Enabled = true
	cfg.Alerts.OpsGenie.APIKey = "test-key-123"
	cfg.Alerts.OpsGenie.Tags = []string{"portwatch", "production"}
	if err := validateOpsGenie(cfg); err != nil {
		t.Errorf("expected no error with tags, got %v", err)
	}
}
