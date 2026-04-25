//go:build integration
// +build integration

package alert_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/alert"
)

// TestSyslogHandler_LocalSocket requires a running syslog daemon and is only
// executed when the "integration" build tag is provided:
//
//	go test -tags integration ./internal/alert/...
func TestSyslogHandler_LocalSocket(t *testing.T) {
	h, err := alert.NewSyslogHandler("", "", "portwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	t.Cleanup(func() { _ = h.Close() })

	events := []alert.Event{
		{Port: 80, Change: "opened", Level: alert.LevelInfo, Time: time.Now()},
		{Port: 443, Change: "closed", Level: alert.LevelAlert, Time: time.Now()},
	}
	for _, e := range events {
		if err := h.Send(e); err != nil {
			t.Errorf("Send(%+v): %v", e, err)
		}
	}
}
