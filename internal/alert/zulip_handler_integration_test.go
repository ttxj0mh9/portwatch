//go:build integration

package alert_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/snapshot"
)

func TestZulipHandler_RealServer(t *testing.T) {
	received := make(chan struct{}, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
		}
		if r.FormValue("type") != "stream" {
			t.Errorf("expected type=stream, got %q", r.FormValue("type"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"success","msg":"","id":1}`))
		received <- struct{}{}
	}))
	defer ts.Close()

	h, err := alert.NewZulipHandler(ts.URL, "bot@example.com", "test-api-key", "general", "portwatch")
	if err != nil {
		t.Fatalf("NewZulipHandler: %v", err)
	}

	ev := alert.NewEvent(snapshot.New([]uint16{443}), "opened")
	if err := h.Send(ev); err != nil {
		t.Fatalf("Send: %v", err)
	}

	select {
	case <-received:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout: server never received request")
	}
}
