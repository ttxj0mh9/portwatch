package config

import (
	"testing"
)

func TestMatrixDefaults_NoOp(t *testing.T) {
	cfg := &Config{}
	matrixDefaults(cfg)
	if cfg.Matrix.Homeserver != "" {
		t.Errorf("expected empty homeserver, got %q", cfg.Matrix.Homeserver)
	}
}

func TestValidateMatrix_Disabled(t *testing.T) {
	cfg := &Config{}
	cfg.Matrix.Enabled = false
	if err := validateMatrix(cfg); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateMatrix_MissingHomeserver(t *testing.T) {
	cfg := &Config{}
	cfg.Matrix.Enabled = true
	cfg.Matrix.AccessToken = "token"
	cfg.Matrix.RoomID = "!room:example.com"
	if err := validateMatrix(cfg); err == nil {
		t.Error("expected error for missing homeserver")
	}
}

func TestValidateMatrix_MissingToken(t *testing.T) {
	cfg := &Config{}
	cfg.Matrix.Enabled = true
	cfg.Matrix.Homeserver = "https://matrix.example.com"
	cfg.Matrix.RoomID = "!room:example.com"
	if err := validateMatrix(cfg); err == nil {
		t.Error("expected error for missing access token")
	}
}

func TestValidateMatrix_MissingRoomID(t *testing.T) {
	cfg := &Config{}
	cfg.Matrix.Enabled = true
	cfg.Matrix.Homeserver = "https://matrix.example.com"
	cfg.Matrix.AccessToken = "token"
	if err := validateMatrix(cfg); err == nil {
		t.Error("expected error for missing room ID")
	}
}

func TestValidateMatrix_Valid(t *testing.T) {
	cfg := &Config{}
	cfg.Matrix.Enabled = true
	cfg.Matrix.Homeserver = "https://matrix.example.com"
	cfg.Matrix.AccessToken = "token"
	cfg.Matrix.RoomID = "!room:example.com"
	if err := validateMatrix(cfg); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}

func TestValidateMatrix_EnabledWithAllFields(t *testing.T) {
	cfg := &Config{}
	cfg.Matrix.Enabled = true
	cfg.Matrix.Homeserver = "https://matrix.org"
	cfg.Matrix.AccessToken = "syt_abc123"
	cfg.Matrix.RoomID = "!abc:matrix.org"
	if err := validateMatrix(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
