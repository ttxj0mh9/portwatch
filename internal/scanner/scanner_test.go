package scanner_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// startTestListener opens a TCP listener on an ephemeral port and returns it.
func startTestListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestTCPScanner_DetectsOpenPort(t *testing.T) {
	ln, port := startTestListener(t)
	defer ln.Close()

	s := &scanner.TCPScanner{
		StartPort: port,
		EndPort:   port,
		Timeout:   200 * time.Millisecond,
	}

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}
	if ports[0].Number != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Number)
	}
	if ports[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", ports[0].Protocol)
	}
}

func TestTCPScanner_NoOpenPorts(t *testing.T) {
	// Use a port range that is very unlikely to have listeners in CI.
	s := scanner.NewTCPScanner(19900, 19910)
	s.Timeout = 100 * time.Millisecond

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan() returned error: %v", err)
	}
	// We can't guarantee zero, but the result must be a valid slice.
	if ports == nil {
		t.Error("expected non-nil slice even when no ports are open")
	}
}

func TestPort_String(t *testing.T) {
	p := scanner.Port{Protocol: "tcp", Number: 8080, Address: "127.0.0.1"}
	want := "127.0.0.1:8080 (tcp)"
	if got := p.String(); got != want {
		t.Errorf("Port.String() = %q, want %q", got, want)
	}
}
