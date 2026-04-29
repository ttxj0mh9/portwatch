package alert

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func twilioConfig(serverURL, from string, to []string) config.TwilioConfig {
	return config.TwilioConfig{
		Enabled:    true,
		AccountSID: "ACTEST123",
		AuthToken:  "token123",
		From:       from,
		To:         to,
	}
}

func TestNewTwilioHandler_MissingAccountSID(t *testing.T) {
	cfg := config.TwilioConfig{AuthToken: "t", From: "+1", To: []string{"+2"}}
	_, err := NewTwilioHandler(cfg)
	if err == nil {
		t.Fatal("expected error for missing account_sid")
	}
}

func TestNewTwilioHandler_MissingAuthToken(t *testing.T) {
	cfg := config.TwilioConfig{AccountSID: "AC", From: "+1", To: []string{"+2"}}
	_, err := NewTwilioHandler(cfg)
	if err == nil {
		t.Fatal("expected error for missing auth_token")
	}
}

func TestNewTwilioHandler_MissingFrom(t *testing.T) {
	cfg := config.TwilioConfig{AccountSID: "AC", AuthToken: "t", To: []string{"+2"}}
	_, err := NewTwilioHandler(cfg)
	if err == nil {
		t.Fatal("expected error for missing from_number")
	}
}

func TestNewTwilioHandler_MissingTo(t *testing.T) {
	cfg := config.TwilioConfig{AccountSID: "AC", AuthToken: "t", From: "+1"}
	_, err := NewTwilioHandler(cfg)
	if err == nil {
		t.Fatal("expected error for missing to_numbers")
	}
}

func TestTwilioHandler_Send_Success(t *testing.T) {
	var received int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received++
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if r.FormValue("To") == "" {
			t.Error("expected To field")
		}
		if r.FormValue("Body") == "" {
			t.Error("expected Body field")
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"sid":"SM123"}`))
	}))
	defer ts.Close()

	cfg := twilioConfig(ts.URL, "+15550001111", []string{"+15559998888", "+15557776666"})
	h, err := NewTwilioHandler(cfg)
	if err != nil {
		t.Fatalf("NewTwilioHandler: %v", err)
	}
	// Rewrite transport so requests hit the test server.
	h.client = ts.Client()
	h.client.Transport = rewriteTransport(ts.URL)

	event := NewEvent(8080, EventOpened)
	if err := h.Send(event); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if received != 2 {
		t.Errorf("expected 2 requests (one per recipient), got %d", received)
	}
}

func TestTwilioHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"invalid number"}`))
	}))
	defer ts.Close()

	cfg := twilioConfig(ts.URL, "+15550001111", []string{"+15559998888"})
	h, err := NewTwilioHandler(cfg)
	if err != nil {
		t.Fatalf("NewTwilioHandler: %v", err)
	}
	h.client = ts.Client()
	h.client.Transport = rewriteTransport(ts.URL)

	event := NewEvent(8080, EventOpened)
	if err := h.Send(event); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}
