package infra

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

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
//
// Important: if using this option, WithPassword and WithAddr won't have any
// effect, because those settings are expected to be part of the client
// definition itself.
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

	if redisCache.client != nil {
		return &redisCache
	}

	if redisCache.addr == "" {
		redisCache.addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisCache.addr,
		Password: redisCache.password,
		DB:       0, // use default DB
	})
	redisCache.client = client

	return &redisCache
}

// RedisCache represents the Redis integration service.
type RedisCache struct {
	client   *redis.Client
	addr     string
	password string
}

// Incr increments the integer in key by 1 with a TTL defined by expiration on Redis.
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

// Decr decrements the integer in key by 1 on Redis.
func (r RedisCache) Decr(ctx context.Context, key string) error {
	count, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis decr: %w", err)
	}

	// if the count reaches 0, delete the key to save cache memory
	if count <= 0 {
		err = r.client.Del(ctx, key).Err()
		if err != nil {
			return fmt.Errorf("redis del: %w", err)
		}
	}

	return nil
}

// Get retrieves the value for the given cache key on Redis.
func (r RedisCache) Get(ctx context.Context, key string) string {
	return r.client.Get(ctx, key).Val()
}

// Set sets a new key/value pair to the Redis cache.
func (r RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Ping checks if Redis connection is healthy.
func (r RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
