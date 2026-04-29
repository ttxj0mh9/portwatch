package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// kafkaWriter is an interface to allow mocking in tests.
type kafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// KafkaHandler publishes port-change events to a Kafka topic.
type KafkaHandler struct {
	writer kafkaWriter
	topic  string
}

// NewKafkaHandler creates a KafkaHandler that writes to the given brokers and topic.
func NewKafkaHandler(brokers []string, topic string) (*KafkaHandler, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("kafka: at least one broker address is required")
	}
	if topic == "" {
		return nil, fmt.Errorf("kafka: topic is required")
	}
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 5 * time.Second,
	}
	return &KafkaHandler{writer: w, topic: topic}, nil
}

type kafkaPayload struct {
	Timestamp string `json:"timestamp"`
	Port      int    `json:"port"`
	Event     string `json:"event"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// Send encodes the event as JSON and publishes it to the configured Kafka topic.
func (h *KafkaHandler) Send(e Event) error {
	p := kafkaPayload{
		Timestamp: e.Time.UTC().Format(time.RFC3339),
		Port:      e.Port,
		Event:     string(e.Kind),
		Level:     string(e.Level),
		Message:   FormatAlert(e),
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("kafka: marshal payload: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return h.writer.WriteMessages(ctx, kafka.Message{Value: body})
}
