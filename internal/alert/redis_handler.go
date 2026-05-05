package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisClient abstracts the Redis client for testability.
type redisClient interface {
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Close() error
}

// RedisHandler publishes alert events to a Redis pub/sub channel.
type RedisHandler struct {
	client  redisClient
	channel string
}

type redisPayload struct {
	Level     string `json:"level"`
	Port      int    `json:"port"`
	Proto     string `json:"proto"`
	Event     string `json:"event"`
	Timestamp string `json:"timestamp"`
}

// NewRedisHandler creates a RedisHandler that publishes to the given channel.
// addr is in "host:port" format.
func NewRedisHandler(addr, password string, db int, channel string) (*RedisHandler, error) {
	if addr == "" {
		return nil, fmt.Errorf("redis: addr is required")
	}
	if channel == "" {
		return nil, fmt.Errorf("redis: channel is required")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisHandler{client: client, channel: channel}, nil
}

// Send publishes the event as a JSON message to the configured Redis channel.
func (h *RedisHandler) Send(e Event) error {
	payload := redisPayload{
		Level:     e.Level.String(),
		Port:      e.Port.Number,
		Proto:     e.Port.Proto,
		Event:     eventLabel(e),
		Timestamp: e.Time.Format(time.RFC3339),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("redis: marshal payload: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.client.Publish(ctx, h.channel, string(data)).Err(); err != nil {
		return fmt.Errorf("redis: publish: %w", err)
	}
	return nil
}

func eventLabel(e Event) string {
	if e.Opened {
		return "opened"
	}
	return "closed"
}
