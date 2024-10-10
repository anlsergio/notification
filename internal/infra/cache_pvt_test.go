package infra

import (
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRedisCache(t *testing.T) {
	t.Run("default settings", func(t *testing.T) {
		redisCache := NewRedisCache()
		assert.NotNil(t, redisCache)

		t.Run("client is set", func(t *testing.T) {
			assert.NotNil(t, redisCache.client)
		})
		t.Run("addr is set", func(t *testing.T) {
			assert.NotEmpty(t, redisCache.addr)
		})
	})

	t.Run("address is provided", func(t *testing.T) {
		redisCache := NewRedisCache(WithAddr("127.0.0.1:5555"))
		assert.NotNil(t, redisCache)

		assert.Equal(t, "127.0.0.1:5555", redisCache.addr)
	})

	t.Run("password is provided", func(t *testing.T) {
		redisCache := NewRedisCache(WithPassword("foo123"))
		assert.NotNil(t, redisCache)

		assert.Equal(t, "foo123", redisCache.password)
	})

	t.Run("custom client is provided", func(t *testing.T) {
		clientMock, _ := redismock.NewClientMock()
		redisCache := NewRedisCache(WithClient(clientMock))
		assert.NotNil(t, redisCache)
	})
}
