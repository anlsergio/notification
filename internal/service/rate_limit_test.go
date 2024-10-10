package service_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"notification/internal/service"
	"notification/mocks"
	"testing"
)

func TestCacheRateLimitHandler_Check(t *testing.T) {
	t.Run("is not rate limited", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("1")

		checker := service.NewCacheRateLimitHandler(cacheSvc)
		ok, err := checker.Check(context.Background(), "123", service.Status)
		require.NoError(t, err)
		require.True(t, ok)
	})
	t.Run("is rate limited", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("100")

		checker := service.NewCacheRateLimitHandler(cacheSvc)
		ok, err := checker.Check(context.Background(), "123", service.Status)
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestCacheRateLimitHandler_IncrementCount(t *testing.T) {
	t.Run("increments count", func(t *testing.T) {
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Incr", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		checker := service.NewCacheRateLimitHandler(cacheSvc)
		require.NoError(t, checker.IncrementCount(context.Background(), "123", service.Status))
	})
}
