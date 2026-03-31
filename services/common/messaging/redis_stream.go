package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

type RedisStream struct {
	Client *redis.Client
}

func NewRedisStream(client *redis.Client) *RedisStream {
	return &RedisStream{Client: client}
}

func (rs *RedisStream) PublishEvent(ctx context.Context, streamName string, event Event) (string, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("failed to marshal event: %w", err)
	}

	res, err := rs.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"payload": string(data),
		},
	}).Result()

	if err != nil {
		return "", fmt.Errorf("failed to add event to stream: %w", err)
	}

	return res, nil
}

func (rs *RedisStream) CreateConsumerGroup(ctx context.Context, streamName, groupName string) error {
	err := rs.Client.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	return nil
}

func (rs *RedisStream) Consume(ctx context.Context, streamName, groupName, consumerName string, handler func(event Event) error) error {
	for {
		streams, err := rs.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Count:    10,
			Block:    0,
		}).Result()

		if err != nil {
			return fmt.Errorf("failed to read from stream: %w", err)
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				var event Event
				payloadStr := message.Values["payload"].(string)
				if err := json.Unmarshal([]byte(payloadStr), &event); err != nil {
					continue // Skip invalid events
				}

				if err := handler(event); err == nil {
					rs.Client.XAck(ctx, streamName, groupName, message.ID)
				}
			}
		}
	}
}
