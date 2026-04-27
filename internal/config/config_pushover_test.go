package config

import (
	"testing"
)

func TestPushoverDefaults_NoOp(t *testing.T) {
	cfg := &Config{}
	pushoverDefaults(cfg)
	// No fields are modified by defaults; just ensure no panic.
}

func TestValidatePushover_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Pushover.Enabled = false
	if err := validatePushover(cfg); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestValidatePushover_MissingUserKey(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Pushover.Enabled = true
	cfg.Alerts.Pushover.APIToken = "token123"
	if err := validatePushover(cfg); err == nil {
		t.Fatal("expected error for missing user key")
	}
}

func TestValidatePushover_MissingAPIToken(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Pushover.Enabled = true
	cfg.Alerts.Pushover.UserKey = "userkey123"
	if err := validatePushover(cfg); err == nil {
		t.Fatal("expected error for missing api token")
	}
}

func TestValidatePushover_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Pushover.Enabled = true
	cfg.Alerts.Pushover.UserKey = "userkey123"
	cfg.Alerts.Pushover.APIToken = "token123"
	if err := validatePushover(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidatePushover_EnabledWithBoth(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.Pushover.Enabled = true
	cfg.Alerts.Pushover.UserKey = "u"
	cfg.Alerts.Pushover.APIToken = "t"
	if err := validatePushover(cfg); err != nil {
		t.Errorf("expected valid config, got: %v", err)
	}
}
