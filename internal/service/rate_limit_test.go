package service_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"notification/internal/domain"
	"notification/internal/service"
	"notification/mocks"
	"testing"
	"time"
)

func TestCacheRateLimitHandler_Check(t *testing.T) {
	rules := domain.RateLimitRules{
		domain.Status: domain.RateLimitRule{
			MaxCount:   2,
			Expiration: time.Minute * 1,
		},
		domain.News: domain.RateLimitRule{
			MaxCount:   1,
			Expiration: time.Hour * 24,
		},
		domain.Marketing: domain.RateLimitRule{
			MaxCount:   3,
			Expiration: time.Hour * 1,
		},
	}

	t.Run("is not rate limited", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("1")

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		ok, err := checker.Check(context.Background(), "123", domain.Status)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("when key is not set should default count to 0", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("")

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		ok, err := checker.Check(context.Background(), "123", domain.Status)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("is rate limited", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("100")

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		ok, err := checker.Check(context.Background(), "123", domain.Status)
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("key is invalid type", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("abc")

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		_, err := checker.Check(context.Background(), "123", domain.Status)
		assert.Error(t, err)
	})
}

func TestCacheRateLimitHandler_IncrementCount(t *testing.T) {
	rules := domain.RateLimitRules{
		domain.Status: domain.RateLimitRule{
			MaxCount:   2,
			Expiration: time.Minute * 1,
		},
	}

	t.Run("increments count", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Incr", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		require.NoError(t, checker.IncrementCount(context.Background(), "123", domain.Status))
	})
}
