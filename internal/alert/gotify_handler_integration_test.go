//go:build integration
// +build integration

package alert

import (
	"os"
	"testing"
)

// TestGotifyHandler_RealServer sends a real notification to a Gotify server.
// Requires environment variables:
//
//	GOTIFY_URL   - e.g. http://gotify.example.com
//	GOTIFY_TOKEN - application token
func TestGotifyHandler_RealServer(t *testing.T) {
	url := os.Getenv("GOTIFY_URL")
	token := os.Getenv("GOTIFY_TOKEN")
	if url == "" || token == "" {
		t.Skip("GOTIFY_URL or GOTIFY_TOKEN not set")
	}

	h, err := NewGotifyHandler(url, token)
	if err != nil {
		t.Fatalf("NewGotifyHandler: %v", err)
	}

	event := NewEvent(8080, KindOpened)
	if err := h.Send(event); err != nil {
		t.Fatalf("Send: %v", err)
	}
}
