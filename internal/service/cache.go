package service

import "context"

// Cache is the abstract representation of the Cache service.
type Cache interface {
	// Set sets a value to the cache for the given key.
	Set(ctx context.Context, key string, value any)
	// Get retrieves the value for the given cache key.
	Get(ctx context.Context, key string) any
}

// NewRedisCache instantiates a new RedisCache instance.
func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

// RedisCache is the concrete implementation of the Redis Cache service.
type RedisCache struct{}

func (r RedisCache) Set(ctx context.Context, key string, value any) {
	//TODO implement me
	panic("implement me")
}

// Get retrieves the value for the given cache key on Redis.
func (r RedisCache) Get(ctx context.Context, key string) any {
	//TODO implement me
	panic("implement me")
}
