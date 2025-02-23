// redis/pubsub.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

// PubSubRepository defines an interface for Redis Pub/Sub operations.
type PubSubRepository interface {
	Publish(ctx context.Context, roomName string, message interface{}) error
	Subscribe(ctx context.Context, roomName string, handler func([]byte))
	Unsubscribe(ctx context.Context, roomName string)
}

// pubSubRepository is a concrete implementation of PubSubRepository.
type pubSubRepository struct {
	client      *goredis.Client
	subscribeMu sync.Mutex
}

// NewPubSubRepository creates a new instance of pubSubRepository.
func NewPubSubRepository(rc *RedisClient) PubSubRepository {
	return &pubSubRepository{
		client: rc.GetRawClient(),
	}
}

// Publish publishes a message to a Redis channel for the specified room.
func (r *pubSubRepository) Publish(ctx context.Context, roomName string, message interface{}) error {
	channel := fmt.Sprintf("room:%s", roomName)
	msgJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v\n", err)
		return err
	}

	if err := r.client.Publish(ctx, channel, msgJSON).Err(); err != nil {
		log.Printf("Failed to publish message to channel %s: %v\n", channel, err)
		return err
	}
	return nil
}

// Subscribe listens for messages on the given room's channel and invokes handler on each message.
func (r *pubSubRepository) Subscribe(ctx context.Context, roomName string, handler func([]byte)) {
	channel := fmt.Sprintf("room:%s", roomName)
	r.subscribeMu.Lock()
	defer r.subscribeMu.Unlock()

	for {
		pubsub := r.client.Subscribe(ctx, channel)
		ch := pubsub.Channel()

		for msg := range ch {
			handler([]byte(msg.Payload))
		}

		// If the subscription ends unexpectedly, wait and try again.
		time.Sleep(2 * time.Second)
		log.Printf("Reconnecting to Redis channel %s...", channel)
	}
}

// Unsubscribe unsubscribes from the specified room channel.
func (r *pubSubRepository) Unsubscribe(ctx context.Context, roomName string) {
	channel := fmt.Sprintf("room:%s", roomName)
	pubsub := r.client.Subscribe(ctx, channel)
	if err := pubsub.Unsubscribe(ctx, channel); err != nil {
		log.Printf("Failed to unsubscribe from channel %s: %v\n", channel, err)
	}
	_ = pubsub.Close()
}
