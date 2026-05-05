package alert

import (
	"encoding/json"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

// fakeMQTTToken satisfies mqtt.Token.
type fakeMQTTToken struct{ err error }

func (t *fakeMQTTToken) Wait() bool                        { return true }
func (t *fakeMQTTToken) WaitTimeout(_ time.Duration) bool  { return true }
func (t *fakeMQTTToken) Done() <-chan struct{}              { ch := make(chan struct{}); close(ch); return ch }
func (t *fakeMQTTToken) Error() error                      { return t.err }

// fakeMQTTClient captures publishes.
type fakeMQTTClient struct {
	published [][]byte
	pubErr    error
}

func (c *fakeMQTTClient) Publish(_ string, _ byte, _ bool, payload interface{}) mqtt.Token {
	if b, ok := payload.([]byte); ok {
		c.published = append(c.published, b)
	}
	return &fakeMQTTToken{err: c.pubErr}
}
func (c *fakeMQTTClient) IsConnected() bool      { return true }
func (c *fakeMQTTClient) Disconnect(_ uint)      {}

func newTestMQTTHandler(client mqttClient) *MQTTHandler {
	return &MQTTHandler{client: client, topic: "portwatch/events"}
}

func TestMQTTHandler_Send_Success(t *testing.T) {
	client := &fakeMQTTClient{}
	h := newTestMQTTHandler(client)

	e := NewEvent(scanner.Port{Number: 8080, Proto: "tcp"}, EventOpened)
	if err := h.Send(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(client.published) != 1 {
		t.Fatalf("expected 1 publish, got %d", len(client.published))
	}

	var p mqttPayload
	if err := json.Unmarshal(client.published[0], &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.Port != 8080 {
		t.Errorf("expected port 8080, got %d", p.Port)
	}
	if p.Event != "opened" {
		t.Errorf("expected event opened, got %s", p.Event)
	}
}

func TestMQTTHandler_Send_PublishError(t *testing.T) {
	client := &fakeMQTTClient{pubErr: errTest}
	h := newTestMQTTHandler(client)

	e := NewEvent(scanner.Port{Number: 443, Proto: "tcp"}, EventClosed)
	if err := h.Send(e); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewMQTTHandler_MissingBrokerURL(t *testing.T) {
	cfg := config.MQTTConfig{Topic: "portwatch/events"}
	if _, err := NewMQTTHandler(cfg); err == nil {
		t.Fatal("expected error for missing broker_url")
	}
}

func TestNewMQTTHandler_MissingTopic(t *testing.T) {
	cfg := config.MQTTConfig{BrokerURL: "tcp://localhost:1883"}
	if _, err := NewMQTTHandler(cfg); err == nil {
		t.Fatal("expected error for missing topic")
	}
}
