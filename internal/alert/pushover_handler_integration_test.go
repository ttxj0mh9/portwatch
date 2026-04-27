package alert

import (
	"os"
	"testing"
	"time"
)

// TestPushoverHandler_RealAPI sends a real Pushover notification.
// Requires PUSHOVER_USER_KEY and PUSHOVER_API_TOKEN environment variables.
// Run with: go test -run TestPushoverHandler_RealAPI -tags integration
func TestPushoverHandler_RealAPI(t *testing.T) {
	userKey := os.Getenv("PUSHOVER_USER_KEY")
	apiToken := os.Getenv("PUSHOVER_API_TOKEN")
	if userKey == "" || apiToken == "" {
		t.Skip("PUSHOVER_USER_KEY and PUSHOVER_API_TOKEN not set; skipping integration test")
	}

	h, err := NewPushoverHandler(userKey, apiToken)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	event := Event{
		Port:   8080,
		Change: "opened",
		Level:  LevelAlert,
		Time:   time.Now(),
	}

	if err := h.Send(event); err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}
