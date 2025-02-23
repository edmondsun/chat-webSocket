// redis/lock.go
package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// DistributedLock demonstrates a simplified Redlock-like approach.
type DistributedLock struct {
	Client     *redis.Client
	Key        string
	Value      string
	Expiration time.Duration
}

// NewDistributedLock creates a new DistributedLock instance.
func NewDistributedLock(client *redis.Client, key string, value string, expiration time.Duration) *DistributedLock {
	return &DistributedLock{
		Client:     client,
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}
}

// Acquire tries to acquire the lock using SET NX.
func (dl *DistributedLock) Acquire(ctx context.Context) (bool, error) {
	ok, err := dl.Client.SetNX(ctx, dl.Key, dl.Value, dl.Expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}
	return ok, nil
}

// Release uses a Lua script to release the lock only if the value matches.
func (dl *DistributedLock) Release(ctx context.Context) error {
	luaScript := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	res, err := dl.Client.Eval(ctx, luaScript, []string{dl.Key}, dl.Value).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	if res.(int64) == 0 {
		log.Printf("Lock for key %s not released: value mismatch.", dl.Key)
	} else {
		log.Printf("Lock for key %s released successfully.", dl.Key)
	}
	return nil
}
