package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPHandler_MissingURL(t *testing.T) {
	_, err := NewHTTPHandler("", "", nil)
	if err == nil {
		t.Fatal("expected error for missing url")
	}
}

func TestNewHTTPHandler_DefaultMethod(t *testing.T) {
	h, err := NewHTTPHandler("http://example.com", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.method != http.MethodPost {
		t.Errorf("expected POST, got %s", h.method)
	}
}

func TestNewHTTPHandler_Valid(t *testing.T) {
	h, err := NewHTTPHandler("http://example.com", "PUT", map[string]string{"X-Token": "abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.method != "PUT" {
		t.Errorf("expected PUT, got %s", h.method)
	}
}

func TestHTTPHandler_Send_Success(t *testing.T) {
	var received httpPayload
	var gotMethod string
	var gotContentType string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewHTTPHandler(ts.URL, "POST", nil)
	e := NewEvent(9200, EventOpened)
	e.At = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected application/json, got %s", gotContentType)
	}
	if received.Port != 9200 {
		t.Errorf("expected port 9200, got %d", received.Port)
	}
	if received.Event != EventOpened.String() {
		t.Errorf("unexpected event kind: %s", received.Event)
	}
}

func TestHTTPHandler_Send_CustomHeaders(t *testing.T) {
	var gotHeader string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Api-Key")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	h, _ := NewHTTPHandler(ts.URL, "", map[string]string{"X-Api-Key": "secret"})
	if err := h.Send(NewEvent(80, EventClosed)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotHeader != "secret" {
		t.Errorf("expected header value 'secret', got %q", gotHeader)
	}
}

func TestHTTPHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h, _ := NewHTTPHandler(ts.URL, "", nil)
	if err := h.Send(NewEvent(443, EventOpened)); err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}

func TestHTTPHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewHTTPHandler("http://127.0.0.1:1", "", nil)
	if err := h.Send(NewEvent(8080, EventOpened)); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
