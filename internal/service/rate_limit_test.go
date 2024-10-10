package service_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"notification/internal/service"
	"notification/mocks"
	"testing"
	"time"
)

func TestCacheRateLimitHandler_Check(t *testing.T) {
	rules := service.RateLimitRules{
		service.Status: service.RateLimitRule{
			MaxCount:   2,
			Expiration: time.Minute * 1,
		},
		service.News: service.RateLimitRule{
			MaxCount:   1,
			Expiration: time.Hour * 24,
		},
		service.Marketing: service.RateLimitRule{
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
		ok, err := checker.Check(context.Background(), "123", service.Status)
		require.NoError(t, err)
		require.True(t, ok)
	})
	t.Run("is rate limited", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("100")

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		ok, err := checker.Check(context.Background(), "123", service.Status)
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestCacheRateLimitHandler_IncrementCount(t *testing.T) {
	rules := service.RateLimitRules{
		service.Status: service.RateLimitRule{
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
		require.NoError(t, checker.IncrementCount(context.Background(), "123", service.Status))
	})
}
