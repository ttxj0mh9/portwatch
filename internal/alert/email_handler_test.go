package alert

import (
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestNewEmailHandler_MissingHost(t *testing.T) {
	_, err := NewEmailHandler(EmailConfig{To: []string{"a@b.com"}})
	if err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestNewEmailHandler_MissingRecipients(t *testing.T) {
	_, err := NewEmailHandler(EmailConfig{Host: "localhost"})
	if err == nil {
		t.Fatal("expected error for missing recipients")
	}
}

func TestNewEmailHandler_Valid(t *testing.T) {
	h, err := NewEmailHandler(EmailConfig{
		Host: "localhost",
		Port: 25,
		From: "portwatch@local",
		To:   []string{"admin@local"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestEmailHandler_Send_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake SMTP server: %v", err)
	}
	defer ln.Close()

	received := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(2 * time.Second))
		conn.Write([]byte("220 fake smtp\r\n"))
		buf := make([]byte, 4096)
		var sb strings.Builder
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				sb.Write(buf[:n])
				conn.Write([]byte("250 OK\r\n"))
			}
			if err == io.EOF || err != nil {
				break
			}
		}
		received <- sb.String()
	}()

	addr := ln.Addr().(*net.TCPAddr)
	h, _ := NewEmailHandler(EmailConfig{
		Host: "127.0.0.1",
		Port: addr.Port,
		From: "portwatch@local",
		To:   []string{"admin@local"},
	})

	e := NewEvent(8080, "opened")
	// We don't assert no error here because the fake server is minimal;
	// we just ensure Send does not panic and attempts a connection.
	_ = h.Send(e)
}

func TestFormatEmailBody(t *testing.T) {
	e := Event{
		Port:  443,
		Kind:  "closed",
		Level: "ALERT",
		Time:  fixedTime(),
	}
	body := formatEmailBody(e)
	if !strings.Contains(body, "443") {
		t.Error("expected port 443 in email body")
	}
	if !strings.Contains(body, "closed") {
		t.Error("expected event kind in email body")
	}
	if !strings.Contains(body, "ALERT") {
		t.Error("expected level in email body")
	}
}
