package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewGrafanaHandler_MissingURL(t *testing.T) {
	_, err := NewGrafanaHandler("", "key")
	if err == nil {
		t.Fatal("expected error for missing baseURL")
	}
}

func TestNewGrafanaHandler_MissingAPIKey(t *testing.T) {
	_, err := NewGrafanaHandler("http://localhost:3000", "")
	if err == nil {
		t.Fatal("expected error for missing apiKey")
	}
}

func TestNewGrafanaHandler_Valid(t *testing.T) {
	h, err := NewGrafanaHandler("http://localhost:3000", "mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestGrafanaHandler_Send_Success(t *testing.T) {
	var received grafanaPayload
	var authHeader string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/annotations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		authHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewGrafanaHandler(ts.URL, "testkey")
	ev := Event{
		Type:   EventOpened,
		Port:   8080,
		Time:   time.Unix(1700000000, 0),
		Level:  LevelAlert,
	}

	if err := h.Send(ev); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if authHeader != "Bearer testkey" {
		t.Errorf("expected Bearer testkey, got %s", authHeader)
	}
	if len(received.Tags) == 0 {
		t.Error("expected tags in payload")
	}
	found := false
	for _, tag := range received.Tags {
		if tag == "opened" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'opened' tag, got %v", received.Tags)
	}
}

func TestGrafanaHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	h, _ := NewGrafanaHandler(ts.URL, "badkey")
	ev := Event{Type: EventClosed, Port: 22, Time: time.Now(), Level: LevelInfo}

	if err := h.Send(ev); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestGrafanaHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewGrafanaHandler("http://127.0.0.1:19999", "key")
	ev := Event{Type: EventOpened, Port: 80, Time: time.Now(), Level: LevelAlert}
	if err := h.Send(ev); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
