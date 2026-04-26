package config

import (
	"testing"
)

func TestValidateSNS_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.SNS.Enabled = false
	if err := validateSNS(cfg); err != nil {
		t.Fatalf("expected no error when SNS disabled, got: %v", err)
	}
}

func TestValidateSNS_MissingTopicARN(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.SNS.Enabled = true
	cfg.Alerts.SNS.TopicARN = ""
	if err := validateSNS(cfg); err == nil {
		t.Fatal("expected error for missing topic_arn")
	}
}

func TestValidateSNS_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.SNS.Enabled = true
	cfg.Alerts.SNS.TopicARN = "arn:aws:sns:us-east-1:123456789012:portwatch-alerts"
	if err := validateSNS(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSNSDefaults_NoOp(t *testing.T) {
	cfg := &Config{}
	// snsDefaults should not panic or mutate meaningful fields
	snsDefaults(cfg)
	if cfg.Alerts.SNS.Enabled {
		t.Error("expected SNS to remain disabled after defaults")
	}
}
