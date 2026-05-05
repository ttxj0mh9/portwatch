//go:build integration
// +build integration

package alert_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestMatrixHandler_RealServer(t *testing.T) {
	received := make(chan struct{}, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"event_id":"$abc123"}`))
		received <- struct{}{}
	}))
	defer server.Close()

	h, err := alert.NewMatrixHandler(server.URL, "fake-token", "!room:example.com")
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	event := alert.NewEvent(scanner.Port{Number: 8080, Proto: "tcp"}, alert.EventOpened)
	if err := h.Send(event); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	select {
	case <-received:
		// success
	case <-time.After(3 * time.Second):
		t.Error("timed out waiting for server to receive a request")
	}
}
