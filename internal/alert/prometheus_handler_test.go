package alert

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not find free port: %v", err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestNewPrometheusHandler_MissingAddr(t *testing.T) {
	_, err := NewPrometheusHandler("", "/metrics")
	if err == nil {
		t.Fatal("expected error for empty addr, got nil")
	}
}

func TestPrometheusHandler_Send_IncrementsCounters(t *testing.T) {
	addr := freeAddr(t)
	h, err := NewPrometheusHandler(addr, "/metrics")
	if err != nil {
		t.Fatalf("NewPrometheusHandler: %v", err)
	}
	defer h.Close()

	// Allow the server a moment to start.
	time.Sleep(50 * time.Millisecond)

	openedEvent := NewEvent(KindOpened, 8080, fixedTime())
	closedEvent := NewEvent(KindClosed, 22, fixedTime())

	for i := 0; i < 3; i++ {
		if err := h.Send(openedEvent); err != nil {
			t.Fatalf("Send opened: %v", err)
		}
	}
	if err := h.Send(closedEvent); err != nil {
		t.Fatalf("Send closed: %v", err)
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	text := string(body)

	if !strings.Contains(text, "portwatch_ports_opened_total 3") {
		t.Errorf("expected opened counter = 3 in metrics output:\n%s", text)
	}
	if !strings.Contains(text, "portwatch_ports_closed_total 1") {
		t.Errorf("expected closed counter = 1 in metrics output:\n%s", text)
	}
}

func TestPrometheusHandler_Send_AlertCounter(t *testing.T) {
	addr := freeAddr(t)
	h, err := NewPrometheusHandler(addr, "/metrics")
	if err != nil {
		t.Fatalf("NewPrometheusHandler: %v", err)
	}
	defer h.Close()
	time.Sleep(50 * time.Millisecond)

	// Port 22 opened => ClassifyEvent should yield LevelAlert.
	e := ClassifyEvent(NewEvent(KindOpened, 22, fixedTime()))
	if err := h.Send(e); err != nil {
		t.Fatalf("Send: %v", err)
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if !strings.Contains(string(body), "portwatch_alerts_total 1") {
		t.Errorf("expected alerts counter = 1:\n%s", string(body))
	}
}
