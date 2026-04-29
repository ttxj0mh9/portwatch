package config

import "testing"

func TestSignalRDefaults_NoOp(t *testing.T) {
	cfg := &Config{}
	signalRDefaults(cfg)
	if cfg.Alerts.SignalR.Enabled {
		t.Error("expected enabled to remain false")
	}
	if cfg.Alerts.SignalR.HubURL != "" {
		t.Error("expected hub_url to remain empty")
	}
}

func TestValidateSignalR_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.SignalR.Enabled = false
	if err := validateSignalR(cfg); err != nil {
		t.Fatalf("expected no error when disabled, got %v", err)
	}
}

func TestValidateSignalR_MissingHubURL(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.SignalR.Enabled = true
	cfg.Alerts.SignalR.AccessKey = "key"
	if err := validateSignalR(cfg); err == nil {
		t.Fatal("expected error for missing hub_url")
	}
}

func TestValidateSignalR_MissingAccessKey(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.SignalR.Enabled = true
	cfg.Alerts.SignalR.HubURL = "https://example.service.signalr.net/api/v1/hubs/portwatch"
	if err := validateSignalR(cfg); err == nil {
		t.Fatal("expected error for missing access_key")
	}
}

func TestValidateSignalR_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Alerts.SignalR.Enabled = true
	cfg.Alerts.SignalR.HubURL = "https://example.service.signalr.net/api/v1/hubs/portwatch"
	cfg.Alerts.SignalR.AccessKey = "supersecret"
	if err := validateSignalR(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
