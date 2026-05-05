package alert

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// fakeRedisClient records the last published message.
type fakeRedisClient struct {
	lastChannel string
	lastMessage string
	err         error
}

func (f *fakeRedisClient) Publish(_ context.Context, channel string, message interface{}) *redis.IntCmd {
	f.lastChannel = channel
	f.lastMessage, _ = message.(string)
	cmd := redis.NewIntCmd(context.Background())
	if f.err != nil {
		cmd.SetErr(f.err)
	}
	return cmd
}

func (f *fakeRedisClient) Close() error { return nil }

func TestNewRedisHandler_MissingAddr(t *testing.T) {
	_, err := NewRedisHandler("", "", 0, "portwatch")
	if err == nil {
		t.Fatal("expected error for missing addr")
	}
}

func TestNewRedisHandler_MissingChannel(t *testing.T) {
	_, err := NewRedisHandler("localhost:6379", "", 0, "")
	if err == nil {
		t.Fatal("expected error for missing channel")
	}
}

func TestNewRedisHandler_Valid(t *testing.T) {
	h, err := NewRedisHandler("localhost:6379", "", 0, "portwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestRedisHandler_Send_Success(t *testing.T) {
	fake := &fakeRedisClient{}
	h := &RedisHandler{client: fake, channel: "portwatch"}
	e := Event{
		Level:  LevelAlert,
		Port:   testPort(8080),
		Opened: true,
		Time:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.lastChannel != "portwatch" {
		t.Errorf("channel = %q, want %q", fake.lastChannel, "portwatch")
	}
	var payload redisPayload
	if err := json.Unmarshal([]byte(fake.lastMessage), &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload.Port != 8080 {
		t.Errorf("port = %d, want 8080", payload.Port)
	}
	if payload.Event != "opened" {
		t.Errorf("event = %q, want \"opened\"", payload.Event)
	}
}

func TestRedisHandler_Send_ClosedEvent(t *testing.T) {
	fake := &fakeRedisClient{}
	h := &RedisHandler{client: fake, channel: "portwatch"}
	e := Event{Level: LevelInfo, Port: testPort(443), Opened: false, Time: time.Now()}
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload redisPayload
	json.Unmarshal([]byte(fake.lastMessage), &payload) //nolint:errcheck
	if payload.Event != "closed" {
		t.Errorf("event = %q, want \"closed\"", payload.Event)
	}
}

func TestRedisHandler_Send_PublishError(t *testing.T) {
	fake := &fakeRedisClient{err: errors.New("connection refused")}
	h := &RedisHandler{client: fake, channel: "portwatch"}
	e := Event{Level: LevelAlert, Port: testPort(22), Opened: true, Time: time.Now()}
	if err := h.Send(e); err == nil {
		t.Fatal("expected error from publish failure")
	}
}
