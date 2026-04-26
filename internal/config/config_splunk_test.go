package config

import (
	"testing"
)

func TestSplunkDefaults_Source(t *testing.T) {
	cfg := &Config{}
	splunkDefaults(cfg)
	if cfg.Alerts.Splunk.Source != "portwatch" {
		t.Errorf("expected default source 'portwatch', got %q", cfg.Alerts.Splunk.Source)
	}
}

func TestSplunkDefaults_DoesNotOverride(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Splunk.Source = "myapp"
	splunkDefaults(cfg)
	if cfg.Alerts.Splunk.Source != "myapp" {
		t.Errorf("expected source 'myapp', got %q", cfg.Alerts.Splunk.Source)
	}
}

func TestValidateSplunk_Disabled(t *testing.T) {
	cfg := &Config{}
	if err := validateSplunk(cfg); err != nil {
		t.Fatalf("expected no error when disabled, got %v", err)
	}
}

func TestValidateSplunk_MissingHECURL(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Splunk.Enabled = true
	cfg.Alerts.Splunk.Token = "tok"
	if err := validateSplunk(cfg); err == nil {
		t.Fatal("expected error for missing hec_url")
	}
}

func TestValidateSplunk_MissingToken(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Splunk.Enabled = true
	cfg.Alerts.Splunk.HECURL = "http://splunk:8088/services/collector"
	if err := validateSplunk(cfg); err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestValidateSplunk_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Splunk.Enabled = true
	cfg.Alerts.Splunk.HECURL = "http://splunk:8088/services/collector"
	cfg.Alerts.Splunk.Token = "mytoken"
	if err := validateSplunk(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
