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
}
