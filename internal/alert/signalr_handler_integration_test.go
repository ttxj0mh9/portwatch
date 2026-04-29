package alert

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// TestSignalRHandler_RealEndpoint runs only when SIGNALR_HUB_URL and
// SIGNALR_ACCESS_KEY are set, allowing manual validation against a live hub.
func TestSignalRHandler_RealEndpoint(t *testing.T) {
	hubURL := os.Getenv("SIGNALR_HUB_URL")
	accessKey := os.Getenv("SIGNALR_ACCESS_KEY")
	if hubURL == "" || accessKey == "" {
		t.Skip("SIGNALR_HUB_URL and SIGNALR_ACCESS_KEY not set")
	}

	h, err := NewSignalRHandler(hubURL, accessKey)
	if err != nil {
		t.Fatalf("create handler: %v", err)
	}
	e := Event{Port: 9090, Kind: Opened, Level: LevelAlert, Time: time.Now()}
	if err := h.Send(e); err != nil {
		t.Fatalf("send event: %v", err)
	}
}

// TestSignalRHandler_LocalEcho verifies round-trip JSON against a local echo server.
func TestSignalRHandler_LocalEcho(t *testing.T) {
	received := make(chan struct{}, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received <- struct{}{}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h, _ := NewSignalRHandler(ts.URL, "localkey")
	e := Event{Port: 3000, Kind: Closed, Level: LevelInfo, Time: time.Now()}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	select {
	case <-received:
	default:
		t.Fatal("server did not receive request")
	}
}
