package alert

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/user/portwatch/internal/config"
)

// mqttClient is an interface for testability.
type mqttClient interface {
	Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token
	IsConnected() bool
	Disconnect(quiesce uint)
}

// MQTTHandler publishes port-change events to an MQTT broker.
type MQTTHandler struct {
	client mqttClient
	topic  string
}

type mqttPayload struct {
	Event     string `json:"event"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
}

// NewMQTTHandler creates a new MQTTHandler from config.
func NewMQTTHandler(cfg config.MQTTConfig) (*MQTTHandler, error) {
	if cfg.BrokerURL == "" {
		return nil, fmt.Errorf("mqtt: broker_url is required")
	}
	if cfg.Topic == "" {
		return nil, fmt.Errorf("mqtt: topic is required")
	}

	opts := mqtt.NewClientOptions().
		AddBroker(cfg.BrokerURL).
		SetClientID(cfg.ClientID).
		SetConnectTimeout(5 * time.Second)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	client := mqtt.NewClient(opts)
	if tok := client.Connect(); tok.Wait() && tok.Error() != nil {
		return nil, fmt.Errorf("mqtt: connect: %w", tok.Error())
	}

	return &MQTTHandler{client: client, topic: cfg.Topic}, nil
}

// Send publishes the event to the configured MQTT topic.
func (h *MQTTHandler) Send(e Event) error {
	p := mqttPayload{
		Event:     e.Type.String(),
		Port:      e.Port.Number,
		Protocol:  e.Port.Proto,
		Timestamp: e.Time.UTC().Format(time.RFC3339),
		Level:     e.Level.String(),
	}

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("mqtt: marshal: %w", err)
	}

	tok := h.client.Publish(h.topic, 1, false, data)
	tok.Wait()
	if tok.Error() != nil {
		return fmt.Errorf("mqtt: publish: %w", tok.Error())
	}
	return nil
}
