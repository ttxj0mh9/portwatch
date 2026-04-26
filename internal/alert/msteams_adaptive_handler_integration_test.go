//go:build integration
// +build integration

package alert

import (
	"os"
	"testing"
	"time"
)

// TestMSTeamsAdaptiveHandler_RealWebhook sends a real alert to a Teams webhook.
// Run with: TEAMS_ADAPTIVE_WEBHOOK_URL=https://... go test -tags integration ./internal/alert/
func TestMSTeamsAdaptiveHandler_RealWebhook(t *testing.T) {
	url := os.Getenv("TEAMS_ADAPTIVE_WEBHOOK_URL")
	if url == "" {
		t.Skip("TEAMS_ADAPTIVE_WEBHOOK_URL not set")
	}

	h, err := NewMSTeamsAdaptiveHandler(url)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	e := Event{
		Port:    9090,
		Message: "[integration test] port 9090 opened unexpectedly",
		Level:   LevelAlert,
		Time:    time.Now(),
	}

	if err := h.Send(e); err != nil {
		t.Errorf("Send failed: %v", err)
	}
}
