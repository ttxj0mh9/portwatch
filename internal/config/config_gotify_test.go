package config

import (
	"testing"
)

func TestGotifyDefaults_NoOp(t *testing.T) {
	cfg := Default()
	if cfg.Gotify.Enabled {
		t.Error("gotify should be disabled by default")
	}
	if cfg.Gotify.ServerURL != "" {
		t.Errorf("expected empty server_url, got %q", cfg.Gotify.ServerURL)
	}
}

func TestValidateGotify_Disabled(t *testing.T) {
	cfg := Default()
	cfg.Gotify.Enabled = false
	if err := validateGotify(cfg); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateGotify_MissingServerURL(t *testing.T) {
	cfg := Default()
	cfg.Gotify.Enabled = true
	cfg.Gotify.Token = "tok"
	if err := validateGotify(cfg); err == nil {
		t.Fatal("expected error for missing server_url")
	}
}

func TestValidateGotify_MissingToken(t *testing.T) {
	cfg := Default()
	cfg.Gotify.Enabled = true
	cfg.Gotify.ServerURL = "http://gotify.example.com"
	if err := validateGotify(cfg); err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestValidateGotify_Valid(t *testing.T) {
	cfg := Default()
	cfg.Gotify.Enabled = true
	cfg.Gotify.ServerURL = "http://gotify.example.com"
	cfg.Gotify.Token = "apptoken"
	if err := validateGotify(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoad_GotifyFullConfig(t *testing.T) {
	path := writeTempConfig(t, `
ports: [80]
interval: 10s
gotify:
  enabled: true
  server_url: "http://gotify.local"
  token: "secret"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if !cfg.Gotify.Enabled {
		t.Error("expected gotify enabled")
	}
	if cfg.Gotify.ServerURL != "http://gotify.local" {
		t.Errorf("unexpected server_url: %s", cfg.Gotify.ServerURL)
	}
	if cfg.Gotify.Token != "secret" {
		t.Errorf("unexpected token: %s", cfg.Gotify.Token)
	}
}
