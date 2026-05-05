//go:build integration
// +build integration

package alert

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestRedisHandler_RealPublish requires a running Redis instance on localhost:6379.
// Run with: go test -tags integration ./internal/alert/...
func TestRedisHandler_RealPublish(t *testing.T) {
	const addr = "localhost:6379"
	const channel = "portwatch:test"

	client := redis.NewClient(&redis.Options{Addr: addr})
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available at %s: %v", addr, err)
	}

	sub := client.Subscribe(ctx, channel)
	defer sub.Close()

	h, err := NewRedisHandler(addr, "", 0, channel)
	if err != nil {
		t.Fatalf("NewRedisHandler: %v", err)
	}

	e := Event{
		Level:  LevelAlert,
		Port:   testPort(9090),
		Opened: true,
		Time:   time.Now(),
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("Send: %v", err)
	}

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		t.Fatalf("ReceiveMessage: %v", err)
	}

	var payload redisPayload
	if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload.Port != 9090 {
		t.Errorf("port = %d, want 9090", payload.Port)
	}
	if payload.Event != "opened" {
		t.Errorf("event = %q, want \"opened\"", payload.Event)
	}
}
