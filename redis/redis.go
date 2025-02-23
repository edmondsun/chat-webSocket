// redis/redis.go
package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// RedisClient wraps a raw *redis.Client.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new RedisClient.
func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection.
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully!")

	return &RedisClient{
		client: rdb,
	}
}

// Close closes the underlying Redis connection.
func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

// GetRawClient returns the underlying *redis.Client.
func (rc *RedisClient) GetRawClient() *redis.Client {
	return rc.client
}
