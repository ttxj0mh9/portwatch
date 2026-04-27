package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func TestNewMattermostHandler_MissingURL(t *testing.T) {
	_, err := NewMattermostHandler("", "bot", "")
	if err == nil {
		t.Fatal("expected error for missing webhook URL")
	}
}

func TestNewMattermostHandler_Valid(t *testing.T) {
	h, err := NewMattermostHandler("https://mattermost.example.com/hooks/abc", "portwatch", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestMattermostHandler_Send_Success(t *testing.T) {
	var received mattermostPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewMattermostHandler(ts.URL, "portwatch", "https://example.com/icon.png")
	evt := NewEvent(scanner.Port(8080), EventOpened)

	if err := h.Send(evt); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.Username != "portwatch" {
		t.Errorf("expected username 'portwatch', got %q", received.Username)
	}
	if received.Text == "" {
		t.Error("expected non-empty text")
	}
}

func TestMattermostHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h, _ := NewMattermostHandler(ts.URL, "", "")
	evt := NewEvent(scanner.Port(9090), EventOpened)

	if err := h.Send(evt); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestMattermostHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewMattermostHandler("http://127.0.0.1:0/hook", "", "")
	evt := NewEvent(scanner.Port(443), EventClosed)

	if err := h.Send(evt); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
