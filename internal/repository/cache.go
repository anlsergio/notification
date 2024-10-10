package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// Cache is the abstract representation of the Cache repository.
type Cache interface {
	// Incr increments the integer in key by 1 with a TTL defined by expiration.
	Incr(ctx context.Context, key string, expiration time.Duration) error
	// Get retrieves the value for the given cache key.
	Get(ctx context.Context, key string) string
}

// RedisCacheOption defines the optional parameters for the RedisCache constructor.
type RedisCacheOption func(r *RedisCache)

// WithAddr sets a Redis connection address.
//
// Defaults to "localhost:6379"
func WithAddr(addr string) RedisCacheOption {
	return func(r *RedisCache) {
		r.addr = addr
	}
}

// WithPassword sets a Redis connection password.
func WithPassword(password string) RedisCacheOption {
	return func(r *RedisCache) {
		r.password = password
	}
}

// WithClient defines a custom Redis client.
//
// Defaults to the default redis.Client.
func WithClient(client *redis.Client) RedisCacheOption {
	return func(r *RedisCache) {
		r.client = client
	}
}

// NewRedisCache instantiates a new RedisCache instance.
func NewRedisCache(opts ...RedisCacheOption) *RedisCache {
	redisCache := RedisCache{}
	for _, opt := range opts {
		opt(&redisCache)
	}

	if redisCache.client == nil {
		if redisCache.addr == "" {
			redisCache.addr = "localhost:6379"
		}

		client := redis.NewClient(&redis.Options{
			Addr:     redisCache.addr,
			Password: redisCache.password,
			DB:       0, // use default DB
		})
		redisCache.client = client
	}

	return &redisCache
}

// RedisCache is the concrete implementation of the Redis Cache repository.
type RedisCache struct {
	client   *redis.Client
	addr     string
	password string
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
