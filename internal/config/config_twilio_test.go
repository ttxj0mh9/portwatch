package config

import (
	"testing"
)

func TestTwilioDefaults_NoOp(t *testing.T) {
	cfg := Default()
	if cfg.Alerts.Twilio.Enabled {
		t.Error("expected Twilio to be disabled by default")
	}
}

func TestValidateTwilio_Disabled(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Twilio.Enabled = false
	if err := validateTwilio(cfg); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateTwilio_MissingAccountSID(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Twilio.Enabled = true
	cfg.Alerts.Twilio.AuthToken = "token"
	cfg.Alerts.Twilio.FromNumber = "+15550001111"
	cfg.Alerts.Twilio.ToNumber = "+15559998888"
	if err := validateTwilio(cfg); err == nil {
		t.Error("expected error for missing account_sid")
	}
}

func TestValidateTwilio_MissingAuthToken(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Twilio.Enabled = true
	cfg.Alerts.Twilio.AccountSID = "ACtest"
	cfg.Alerts.Twilio.FromNumber = "+15550001111"
	cfg.Alerts.Twilio.ToNumber = "+15559998888"
	if err := validateTwilio(cfg); err == nil {
		t.Error("expected error for missing auth_token")
	}
}

func TestValidateTwilio_MissingFromNumber(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Twilio.Enabled = true
	cfg.Alerts.Twilio.AccountSID = "ACtest"
	cfg.Alerts.Twilio.AuthToken = "token"
	cfg.Alerts.Twilio.ToNumber = "+15559998888"
	if err := validateTwilio(cfg); err == nil {
		t.Error("expected error for missing from_number")
	}
}

func TestValidateTwilio_MissingToNumber(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Twilio.Enabled = true
	cfg.Alerts.Twilio.AccountSID = "ACtest"
	cfg.Alerts.Twilio.AuthToken = "token"
	cfg.Alerts.Twilio.FromNumber = "+15550001111"
	if err := validateTwilio(cfg); err == nil {
		t.Error("expected error for missing to_number")
	}
}

func TestValidateTwilio_Valid(t *testing.T) {
	cfg := Default()
	cfg.Alerts.Twilio.Enabled = true
	cfg.Alerts.Twilio.AccountSID = "ACtest"
	cfg.Alerts.Twilio.AuthToken = "token"
	cfg.Alerts.Twilio.FromNumber = "+15550001111"
	cfg.Alerts.Twilio.ToNumber = "+15559998888"
	if err := validateTwilio(cfg); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}
