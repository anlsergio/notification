package service_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"notification/internal/service"
	"testing"
)

type stubCache struct {
	notificationCounts int
}

func (s stubCache) Set(ctx context.Context, key string, value any) {}

func (s stubCache) Get(ctx context.Context, key string) any {
	return nil
}

func TestEmailNotification_Send(t *testing.T) {
	t.Run("notification is sent", func(t *testing.T) {
		svc := service.NewEmailNotification(stubCache{1})
		assert.NoError(t, svc.Send("user1", "hey!", service.Marketing))
	})
}
