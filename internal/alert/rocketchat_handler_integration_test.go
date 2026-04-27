package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/snapshot"
)

// TestRocketChatHandler_RealWebhook exercises the full HTTP round-trip against a
// local test server that mimics a Rocket.Chat incoming webhook endpoint.
func TestRocketChatHandler_RealWebhook(t *testing.T) {
	var receivedText string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			t.Errorf("unexpected content-type: %s", ct)
		}
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		receivedText, _ = payload["text"].(string)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, err := alert.NewRocketChatHandler(ts.URL, "#security", "portwatch-bot")
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	snap := snapshot.New([]uint16{443})
	event := alert.NewEvent(snap, alert.EventOpened)

	if err := h.Send(event, snap); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if receivedText == "" {
		t.Error("expected non-empty text in webhook payload")
	}
}
