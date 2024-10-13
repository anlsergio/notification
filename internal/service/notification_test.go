package service_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"notification/internal/domain"
	"notification/internal/repository"
	"notification/internal/service"
	"notification/mocks"
	"testing"
	"time"
)

func TestEmailNotification_Send(t *testing.T) {
	t.Run("notification is sent", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("IsRateLimited", mock.Anything, mock.Anything, mock.Anything).
			Return(true, time.Duration(0), nil)

		mailer := mocks.NewMailSender(t)
		mailer.
			On("SendEmail", mock.Anything, mock.Anything).
			Return(nil)

		userRepo := mocks.NewUserRepository(t)
		userRepo.
			On("Get", mock.Anything).
			Return(domain.User{}, nil)

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo)
		_, err := svc.Send(context.Background(), "user1", "hey!", domain.Marketing)
		assert.NoError(t, err)
	})

	t.Run("rate limit exceeded", func(t *testing.T) {
		retryAfter := time.Minute

		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("IsRateLimited", mock.Anything, mock.Anything, mock.Anything).
			Return(false, retryAfter, nil).
			Maybe()

		mailer := mocks.NewMailSender(t)
		mailer.
			On("SendEmail", mock.Anything, mock.Anything).
			Return(nil).
			Maybe()

		userRepo := mocks.NewUserRepository(t)
		userRepo.
			On("Get", mock.Anything).
			Return(domain.User{}, nil)

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo)
		gotRetryAfter, err := svc.Send(context.Background(), "user1", "hey!", domain.Marketing)
		assert.ErrorIs(t, err, service.ErrRateLimitExceeded)
		assert.Equal(t, retryAfter, gotRetryAfter)

		mailer.AssertNotCalled(t, "SendEmail")
	})

	t.Run("invalid user", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("IsRateLimited", mock.Anything, mock.Anything, mock.Anything).
			Return(false, time.Duration(0), nil).
			Maybe()

		mailer := mocks.NewMailSender(t)
		mailer.
			On("SendEmail", mock.Anything, mock.Anything).
			Return(nil).
			Maybe()

		userRepo := mocks.NewUserRepository(t)
		userRepo.
			On("Get", mock.Anything).
			Return(domain.User{}, repository.ErrInvalidUserID)

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo)
		_, err := svc.Send(context.Background(), "user1", "hey!", domain.Marketing)
		assert.Error(t, err)

		rateLimitHandler.AssertNotCalled(t, "IsRateLimited")
		mailer.AssertNotCalled(t, "SendEmail")
	})
}
