package alert

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

// TestMQTTHandler_RealBroker publishes a real event to a running MQTT broker.
// Set PORTWATCH_TEST_MQTT_BROKER to a broker URL (e.g. tcp://localhost:1883)
// to enable this test.
func TestMQTTHandler_RealBroker(t *testing.T) {
	broker := os.Getenv("PORTWATCH_TEST_MQTT_BROKER")
	if broker == "" {
		t.Skip("PORTWATCH_TEST_MQTT_BROKER not set")
	}

	cfg := config.MQTTConfig{
		Enabled:   true,
		BrokerURL: broker,
		Topic:     "portwatch/test",
		ClientID:  "portwatch-integration-test",
	}

	h, err := NewMQTTHandler(cfg)
	if err != nil {
		t.Fatalf("NewMQTTHandler: %v", err)
	}
	defer h.client.Disconnect(250)

	e := NewEvent(scanner.Port{Number: 9090, Proto: "tcp"}, EventOpened)
	if err := h.Send(e); err != nil {
		t.Fatalf("Send: %v", err)
	}

	// Minimal smoke-check: verify JSON round-trip of the payload shape.
	p := mqttPayload{
		Event:     e.Type.String(),
		Port:      e.Port.Number,
		Protocol:  e.Port.Proto,
		Timestamp: e.Time.UTC().Format(time.RFC3339),
		Level:     e.Level.String(),
	}
	data, _ := json.Marshal(p)
	var back mqttPayload
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("round-trip: %v", err)
	}
	if back.Port != 9090 {
		t.Errorf("expected port 9090, got %d", back.Port)
	}
}
