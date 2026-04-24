package monitor_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func startListener(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	return addr.Port, func() { ln.Close() }
}

func TestMonitor_DetectsOpenedPort(t *testing.T) {
	s := scanner.NewTCPScanner("127.0.0.1", 200*time.Millisecond)

	// Start with no open ports on a chosen port.
	port, stop := startListener(t)
	stop() // close immediately so first scan sees it closed

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	m := monitor.New(s, []int{port}, 150*time.Millisecond)
	go m.Run(ctx)

	// Open the port after the monitor has started.
	time.Sleep(200 * time.Millisecond)
	_, stop2 := startListenerOnPort(t, port)
	defer stop2()

	select {
	case alert := <-m.Alerts:
		if len(alert.Opened) == 0 {
			t.Errorf("expected opened ports, got none")
		}
	case <-ctx.Done():
		t.Error("timed out waiting for alert")
	}
}

func startListenerOnPort(t *testing.T, port int) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", itoa(port)))
	if err != nil {
		t.Fatalf("failed to listen on port %d: %v", port, err)
	}
	return port, func() { ln.Close() }
}

func itoa(n int) string {
	return net.JoinHostPort("", string(rune('0'+n%10)))[1:] // simple, use fmt in real code
}

func TestMonitor_NoAlertWhenNoChange(t *testing.T) {
	s := scanner.NewTCPScanner("127.0.0.1", 200*time.Millisecond)
	port, stop := startListener(t)
	defer stop()

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	m := monitor.New(s, []int{port}, 150*time.Millisecond)
	go m.Run(ctx)

	<-ctx.Done()
	select {
	case alert, ok := <-m.Alerts:
		if ok {
			t.Errorf("unexpected alert: %+v", alert)
		}
	default:
	}
}
