package service_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"notification/internal/service"
	"notification/mocks"
	"testing"
)

func TestEmailNotification_Send(t *testing.T) {
	t.Run("notification is sent", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("Check", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)

		rateLimitHandler.
			On("IncrementCount", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		mailer := mocks.NewMailSender(t)
		mailer.
			On("Send").
			Return(nil)

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer)
		assert.NoError(t, svc.Send(context.Background(), "user1", "hey!", service.Marketing))
	})

	t.Run("rate limit exceeded", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("Check", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil)

		rateLimitHandler.
			On("IncrementCount", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Maybe()

		mailer := mocks.NewMailSender(t)
		mailer.
			On("Send").
			Return(nil).
			Maybe()

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer)
		assert.Error(t, svc.Send(context.Background(), "user1", "hey!", service.Marketing))

		mailer.AssertNotCalled(t, "Send")
		rateLimitHandler.AssertNotCalled(t, "IncrementCount", mock.Anything, mock.Anything, mock.Anything)
	})
}
