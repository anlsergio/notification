package service

import (
	"context"
	"time"
)

// Cache is the abstract representation of the Cache service.
type Cache interface {
	// Incr increments the integer in key by 1 with a TTL defined by expiration.
	Incr(ctx context.Context, key string, expiration time.Duration) error
	// Get retrieves the value for the given cache key.
	Get(ctx context.Context, key string) string
	// Set sets a new key/value pair to the Redis cache.
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	// Decr decrements the integer in key by 1 on Redis.
	Decr(ctx context.Context, key string) error
}
