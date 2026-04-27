package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

// TestMattermostHandler_RealWebhook performs a full round-trip against a local
// httptest server, verifying the complete JSON payload structure.
func TestMattermostHandler_RealWebhook(t *testing.T) {
	var payload mattermostPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("failed to decode payload: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, err := NewMattermostHandler(ts.URL, "portwatch", "https://example.com/icon.png")
	if err != nil {
		t.Fatalf("NewMattermostHandler: %v", err)
	}

	evt := NewEvent(scanner.Port(22), EventOpened)
	if err := h.Send(evt); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if !strings.Contains(payload.Text, "portwatch") {
		t.Errorf("payload text missing '[portwatch]' prefix, got: %q", payload.Text)
	}
	if payload.Username != "portwatch" {
		t.Errorf("expected username 'portwatch', got %q", payload.Username)
	}
	if payload.IconURL != "https://example.com/icon.png" {
		t.Errorf("expected icon URL, got %q", payload.IconURL)
	}
}
