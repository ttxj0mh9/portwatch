//go:build integration
// +build integration

package alert

import (
	"os"
	"testing"
	"time"
)

// TestSplunkHandler_RealHEC sends a real event to a Splunk HEC endpoint.
// Requires environment variables:
//
//	SPLUNK_HEC_URL  — e.g. https://splunk.example.com:8088/services/collector
//	SPLUNK_TOKEN    — HEC token
func TestSplunkHandler_RealHEC(t *testing.T) {
	url := os.Getenv("SPLUNK_HEC_URL")
	token := os.Getenv("SPLUNK_TOKEN")
	if url == "" || token == "" {
		t.Skip("SPLUNK_HEC_URL or SPLUNK_TOKEN not set")
	}

	h, err := NewSplunkHandler(url, token, "portwatch-integration-test")
	if err != nil {
		t.Fatalf("NewSplunkHandler: %v", err)
	}

	e := Event{
		Port:   9200,
		Action: "opened",
		Level:  LevelAlert,
		Time:   time.Now(),
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("Send: %v", err)
	}
}
