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
		cacheSvc.
			On("Incr", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		ok, err := checker.IsRateLimited(context.Background(), "123", domain.Status)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("when key is not set should default count to 0", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("")
		cacheSvc.
			On("Incr", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		ok, err := checker.IsRateLimited(context.Background(), "123", domain.Status)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("is rate limited", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("100")
		cacheSvc.
			On("Incr", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		ok, err := checker.IsRateLimited(context.Background(), "123", domain.Status)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("when check fails it doesn't increment count", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("abc") // the string "abc" causes the string to int conversion to fail.
		cacheSvc.
			On("Incr", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Maybe()

		checker := service.NewCacheRateLimitHandler(cacheSvc, rules)
		_, err := checker.IsRateLimited(context.Background(), "123", domain.Status)
		assert.Error(t, err)
		cacheSvc.AssertNotCalled(t, "Incr", mock.Anything, mock.Anything, mock.Anything)
	})
}
