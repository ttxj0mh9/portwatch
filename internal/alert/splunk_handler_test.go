package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSplunkHandler_MissingURL(t *testing.T) {
	_, err := NewSplunkHandler("", "token", "portwatch")
	if err == nil {
		t.Fatal("expected error for missing HEC URL")
	}
}

func TestNewSplunkHandler_MissingToken(t *testing.T) {
	_, err := NewSplunkHandler("http://splunk:8088/services/collector", "", "portwatch")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewSplunkHandler_DefaultSource(t *testing.T) {
	h, err := NewSplunkHandler("http://splunk:8088/services/collector", "token", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.source != "portwatch" {
		t.Errorf("expected default source 'portwatch', got %q", h.source)
	}
}

func TestSplunkHandler_Send_Success(t *testing.T) {
	var received splunkEvent

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Splunk mytoken" {
			t.Errorf("missing or wrong Authorization header: %s", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	h, _ := NewSplunkHandler(srv.URL, "mytoken", "portwatch")
	e := Event{Port: 9200, Action: "opened", Level: LevelAlert, Time: time.Unix(1700000000, 0)}

	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.Source != "portwatch" {
		t.Errorf("expected source 'portwatch', got %q", received.Source)
	}
	if int(received.Event["port"].(float64)) != 9200 {
		t.Errorf("unexpected port in payload: %v", received.Event["port"])
	}
}

func TestSplunkHandler_Send_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	h, _ := NewSplunkHandler(srv.URL, "badtoken", "portwatch")
	err := h.Send(Event{Port: 80, Action: "opened", Level: LevelAlert, Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestSplunkHandler_Send_UnreachableURL(t *testing.T) {
	h, _ := NewSplunkHandler("http://127.0.0.1:1", "token", "portwatch")
	err := h.Send(Event{Port: 443, Action: "closed", Level: LevelInfo, Time: time.Now()})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
