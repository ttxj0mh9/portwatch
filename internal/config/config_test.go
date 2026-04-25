package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTempConfig(t, "{}\n")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.ScanInterval)
	}
	if cfg.Snapshot.BackupCount != 3 {
		t.Errorf("expected backup_count 3, got %d", cfg.Snapshot.BackupCount)
	}
	if cfg.Alerts.Log == nil {
		t.Error("expected default log handler to be enabled")
	}
}

func TestLoad_FullConfig(t *testing.T) {
	content := `
scan_interval: 10s
ports:
  include: [80, 443]
  exclude: [22]
snapshot:
  path: /tmp/snap.json
  backup_count: 5
alerts:
  teams:
    webhook_url: https://example.webhook.office.com/abc
  slack:
    webhook_url: https://hooks.slack.com/xyz
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.ScanInterval)
	}
	if len(cfg.Ports.Include) != 2 {
		t.Errorf("expected 2 included ports, got %d", len(cfg.Ports.Include))
	}
	if cfg.Alerts.Teams == nil || cfg.Alerts.Teams.WebhookURL == "" {
		t.Error("expected teams config to be populated")
	}
	if cfg.Snapshot.BackupCount != 5 {
		t.Errorf("expected 5, got %d", cfg.Snapshot.BackupCount)
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	path := writeTempConfig(t, "scan_interval: not-a-duration\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	path := writeTempConfig(t, "ports:\n  include: [99999]\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_TeamsConfig(t *testing.T) {
	content := `
alerts:
  teams:
    webhook_url: https://example.webhook.office.com/teams/hook
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Alerts.Teams == nil {
		t.Fatal("expected teams config, got nil")
	}
	if cfg.Alerts.Teams.WebhookURL != "https://example.webhook.office.com/teams/hook" {
		t.Errorf("unexpected teams webhook URL: %s", cfg.Alerts.Teams.WebhookURL)
	}
}
