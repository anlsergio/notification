package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// Cache is the abstract representation of the Cache service.
type Cache interface {
	// Incr increments the integer in key by 1 with a TTL defined by expiration.
	Incr(ctx context.Context, key string, expiration time.Duration) error
	// Get retrieves the value for the given cache key.
	Get(ctx context.Context, key string) string
}

// NewRedisCache instantiates a new RedisCache instance.
func NewRedisCache() *RedisCache {
	// TODO: inject configs properly
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisCache{client}
}

// RedisCache is the concrete implementation of the Redis Cache service.
type RedisCache struct {
	client *redis.Client
}

func (r RedisCache) Incr(ctx context.Context, key string, expiration time.Duration) error {
	_, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis incr: %w", err)
	}

	err = r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis expire: %w", err)
	}

	return nil
}

// Get retrieves the value for the given cache key on Redis.
func (r RedisCache) Get(ctx context.Context, key string) string {
	return r.client.Get(ctx, key).Val()
}
