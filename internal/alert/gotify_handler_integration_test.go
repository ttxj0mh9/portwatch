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

// TestGotifyHandler_RealServer_ClosedEvent sends a port-closed notification to
// a Gotify server, verifying that both event kinds can be delivered successfully.
// Requires the same environment variables as TestGotifyHandler_RealServer.
func TestGotifyHandler_RealServer_ClosedEvent(t *testing.T) {
	url := os.Getenv("GOTIFY_URL")
	token := os.Getenv("GOTIFY_TOKEN")
	if url == "" || token == "" {
		t.Skip("GOTIFY_URL or GOTIFY_TOKEN not set")
	}

	h, err := NewGotifyHandler(url, token)
	if err != nil {
		t.Fatalf("NewGotifyHandler: %v", err)
	}

	event := NewEvent(8080, KindClosed)
	if err := h.Send(event); err != nil {
		t.Fatalf("Send: %v", err)
	}
}
