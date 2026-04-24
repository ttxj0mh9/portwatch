package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewOpsGenieHandler_MissingKey(t *testing.T) {
	_, err := NewOpsGenieHandler("")
	if err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestNewOpsGenieHandler_Valid(t *testing.T) {
	h, err := NewOpsGenieHandler("test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestOpsGenieHandler_Send_Success(t *testing.T) {
	var received opsGeniePayload
	var authHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	h, _ := NewOpsGenieHandler("my-secret-key")
	h.client = server.Client()
	h.client.Transport = rewriteTransport(server.URL)

	e := NewEvent(Port{Number: 8080, Proto: "tcp"}, true)
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if authHeader != "GenieKey my-secret-key" {
		t.Errorf("expected auth header 'GenieKey my-secret-key', got %q", authHeader)
	}
	if received.Message == "" {
		t.Error("expected non-empty message in payload")
	}
}

func TestOpsGenieHandler_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	h, _ := NewOpsGenieHandler("bad-key")
	h.client = server.Client()
	h.client.Transport = rewriteTransport(server.URL)

	e := NewEvent(Port{Number: 443, Proto: "tcp"}, false)
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestOpsGenieHandler_Send_AlertPriority(t *testing.T) {
	var received opsGeniePayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	h, _ := NewOpsGenieHandler("key")
	h.client = server.Client()
	h.client.Transport = rewriteTransport(server.URL)

	e := NewEvent(Port{Number: 9999, Proto: "tcp"}, true)
	e.Level = LevelAlert
	_ = h.Send(e)

	if received.Priority != "P1" {
		t.Errorf("expected P1 priority for alert level, got %q", received.Priority)
	}
}
