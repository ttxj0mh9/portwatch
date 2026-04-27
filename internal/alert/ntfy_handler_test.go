package alert

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewNtfyHandler_MissingServerURL(t *testing.T) {
	_, err := NewNtfyHandler("", "portwatch")
	if err == nil {
		t.Fatal("expected error for missing server URL")
	}
}

func TestNewNtfyHandler_MissingTopic(t *testing.T) {
	_, err := NewNtfyHandler("https://ntfy.sh", "")
	if err == nil {
		t.Fatal("expected error for missing topic")
	}
}

func TestNewNtfyHandler_Valid(t *testing.T) {
	h, err := NewNtfyHandler("https://ntfy.sh", "portwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestNtfyHandler_Send_Success(t *testing.T) {
	var capturedTitle, capturedPriority string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedTitle = r.Header.Get("Title")
		capturedPriority = r.Header.Get("Priority")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewNtfyHandler(ts.URL, "portwatch")
	h.client = ts.Client()

	e := NewEvent(8080, "opened")
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedTitle == "" {
		t.Error("expected Title header to be set")
	}
	_ = capturedPriority
}

func TestNtfyHandler_Send_AlertPriority(t *testing.T) {
	var capturedPriority, capturedTags string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPriority = r.Header.Get("Priority")
		capturedTags = r.Header.Get("Tags")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, _ := NewNtfyHandler(ts.URL, "portwatch")
	h.client = ts.Client()

	// Port in the watch list triggers LevelAlert
	e := NewEvent(22, "opened")
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedPriority != "high" {
		t.Errorf("expected priority=high, got %q", capturedPriority)
	}
	if capturedTags == "" {
		t.Error("expected Tags header to be set")
	}
}

func TestNtfyHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h, _ := NewNtfyHandler(ts.URL, "portwatch")
	h.client = ts.Client()

	e := NewEvent(8080, "opened")
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestNtfyHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewNtfyHandler("http://127.0.0.1:19999", "portwatch")
	e := NewEvent(8080, "opened")
	if err := h.Send(e); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
