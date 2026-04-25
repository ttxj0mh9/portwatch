package config

import (
	"testing"
)

func TestVictorOpsDefaults_RoutingKey(t *testing.T) {
	cfg := &Config{}
	victorOpsDefaults(cfg)
	if cfg.Alerts.VictorOps.RoutingKey != "default" {
		t.Errorf("expected routing_key 'default', got %q", cfg.Alerts.VictorOps.RoutingKey)
	}
}

func TestVictorOpsDefaults_DoesNotOverride(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.VictorOps.RoutingKey = "my-team"
	victorOpsDefaults(cfg)
	if cfg.Alerts.VictorOps.RoutingKey != "my-team" {
		t.Errorf("expected routing_key 'my-team', got %q", cfg.Alerts.VictorOps.RoutingKey)
	}
}

func TestValidateVictorOps_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.VictorOps.Enabled = false
	if err := validateVictorOps(cfg); err != nil {
		t.Fatalf("expected no error when disabled, got %v", err)
	}
}

func TestValidateVictorOps_MissingRESTURL(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.VictorOps.Enabled = true
	cfg.Alerts.VictorOps.RoutingKey = "default"
	if err := validateVictorOps(cfg); err == nil {
		t.Fatal("expected error for missing rest_url")
	}
}

func TestValidateVictorOps_MissingRoutingKey(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.VictorOps.Enabled = true
	cfg.Alerts.VictorOps.RESTURL = "https://alert.victorops.com/integrations/generic/123/alert/abc"
	cfg.Alerts.VictorOps.RoutingKey = ""
	if err := validateVictorOps(cfg); err == nil {
		t.Fatal("expected error for missing routing_key")
	}
}

func TestValidateVictorOps_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.VictorOps.Enabled = true
	cfg.Alerts.VictorOps.RESTURL = "https://alert.victorops.com/integrations/generic/123/alert/abc"
	cfg.Alerts.VictorOps.RoutingKey = "default"
	if err := validateVictorOps(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
