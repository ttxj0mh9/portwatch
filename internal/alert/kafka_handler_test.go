package alert

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

// mockKafkaWriter captures messages written to it.
type mockKafkaWriter struct {
	messages []kafka.Message
	err      error
}

func (m *mockKafkaWriter) WriteMessages(_ context.Context, msgs ...kafka.Message) error {
	if m.err != nil {
		return m.err
	}
	m.messages = append(m.messages, msgs...)
	return nil
}

func (m *mockKafkaWriter) Close() error { return nil }

func TestNewKafkaHandler_MissingBrokers(t *testing.T) {
	_, err := NewKafkaHandler(nil, "portwatch")
	if err == nil {
		t.Fatal("expected error for missing brokers")
	}
}

func TestNewKafkaHandler_MissingTopic(t *testing.T) {
	_, err := NewKafkaHandler([]string{"localhost:9092"}, "")
	if err == nil {
		t.Fatal("expected error for missing topic")
	}
}

func TestNewKafkaHandler_Valid(t *testing.T) {
	h, err := NewKafkaHandler([]string{"localhost:9092"}, "portwatch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestKafkaHandler_Send_Success(t *testing.T) {
	mock := &mockKafkaWriter{}
	h := &KafkaHandler{writer: mock, topic: "portwatch"}

	e := NewEvent(8080, EventOpened, fixedTime())
	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if len(mock.messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(mock.messages))
	}
	var p struct {
		Port  int    `json:"port"`
		Event string `json:"event"`
		Level string `json:"level"`
	}
	if err := json.Unmarshal(mock.messages[0].Value, &p); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if p.Port != 8080 {
		t.Errorf("expected port 8080, got %d", p.Port)
	}
	if p.Event != string(EventOpened) {
		t.Errorf("expected event %q, got %q", EventOpened, p.Event)
	}
}

func TestKafkaHandler_Send_ClosedEvent(t *testing.T) {
	mock := &mockKafkaWriter{}
	h := &KafkaHandler{writer: mock, topic: "portwatch"}

	e := NewEvent(443, EventClosed, time.Now())
	e = ClassifyEvent(e)
	if err := h.Send(e); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	var p struct {
		Level string `json:"level"`
	}
	if err := json.Unmarshal(mock.messages[0].Value, &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Level != string(LevelAlert) {
		t.Errorf("expected level %q, got %q", LevelAlert, p.Level)
	}
}
