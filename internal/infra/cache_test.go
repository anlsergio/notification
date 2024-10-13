package infra_test

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"notification/internal/infra"
	"testing"
	"time"
)

func TestRedisCache_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	mock.ExpectGet("foo").SetVal("bar")

	redisCache := infra.NewRedisCache(infra.WithClient(db))
	got := redisCache.Get(context.Background(), "foo")
	assert.Equal(t, "bar", got)

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("there were unfulfilled expectations: %v", mock.ExpectationsWereMet())
	}
}

func TestRedisCache_Incr(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	mock.ExpectIncr("foo").SetVal(1)
	mock.ExpectExpire("foo", time.Hour).SetVal(true)

	redisCache := infra.NewRedisCache(infra.WithClient(db))
	err := redisCache.Incr(context.Background(), "foo", time.Hour)
	require.NoError(t, err)

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("there were unfulfilled expectations: %v", mock.ExpectationsWereMet())
	}
}

func TestRedisCache_Decr(t *testing.T) {
	t.Run("decrement counter", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		redisCache := infra.NewRedisCache(infra.WithClient(db))

		key := "foo"
		// initial value before decrement
		expectedValue := int64(5)

		// After decrement, it should be 4.
		mock.ExpectDecr(key).SetVal(expectedValue - 1)

		err := redisCache.Decr(context.Background(), key)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("decrement and delete", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		redisCache := infra.NewRedisCache(infra.WithClient(db))

		// Define the key and expected value.
		key := "foo"

		// After decrement, it should be 0.
		mock.ExpectDecr(key).SetVal(0)
		// The key should be deleted.
		mock.ExpectDel(key).SetVal(1)

		// Call the Decr method.
		err := redisCache.Decr(context.Background(), key)
		assert.NoError(t, err)

		// Verify that the expectations were met.
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
