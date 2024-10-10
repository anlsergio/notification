package repository_test

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"notification/internal/repository"
	"testing"
	"time"
)

func TestRedisCache_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	mock.ExpectGet("foo").SetVal("bar")

	redisCache := repository.NewRedisCache(repository.WithClient(db))
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

	redisCache := repository.NewRedisCache(repository.WithClient(db))
	err := redisCache.Incr(context.Background(), "foo", time.Hour)
	require.NoError(t, err)

	if mock.ExpectationsWereMet() != nil {
		t.Errorf("there were unfulfilled expectations: %v", mock.ExpectationsWereMet())
	}
}
