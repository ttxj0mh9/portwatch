package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
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
	path := writeTempConfig(t, "scan_interval: 60s\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 60*time.Second {
		t.Errorf("expected 60s, got %s", cfg.ScanInterval)
	}
	if !cfg.AlertOnNew {
		t.Error("expected AlertOnNew to default to true")
	}
}

func TestLoad_FullConfig(t *testing.T) {
	content := `
scan_interval: 10s
ports: [22, 80, 443]
alert_on_new: true
alert_on_closed: false
log_file: /var/log/portwatch.log
`
	path := writeTempConfig(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(cfg.Ports))
	}
	if cfg.AlertOnClosed {
		t.Error("expected AlertOnClosed to be false")
	}
	if cfg.LogFile != "/var/log/portwatch.log" {
		t.Errorf("unexpected log_file: %s", cfg.LogFile)
	}
}

func TestLoad_InvalidInterval(t *testing.T) {
	path := writeTempConfig(t, "scan_interval: 500ms\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for short interval")
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	path := writeTempConfig(t, "ports: [0, 80]\n")
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for port 0")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %s", cfg.ScanInterval)
	}
}
