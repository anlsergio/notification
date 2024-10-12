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
)

func TestEmailNotification_Send(t *testing.T) {
	t.Run("notification is sent", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("IsRateLimited", mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil)

		mailer := mocks.NewMailSender(t)
		mailer.
			On("SendEmail", mock.Anything, mock.Anything).
			Return(nil)

		userRepo := mocks.NewUserRepository(t)
		userRepo.
			On("Get", mock.Anything).
			Return(domain.User{}, nil)

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo)
		assert.NoError(t, svc.Send(context.Background(), "user1", "hey!", domain.Marketing))
	})

	t.Run("rate limit exceeded", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("IsRateLimited", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
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
		assert.Error(t, svc.Send(context.Background(), "user1", "hey!", domain.Marketing))

		mailer.AssertNotCalled(t, "SendEmail")
	})

	t.Run("invalid user", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("IsRateLimited", mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).
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
		assert.Error(t, svc.Send(context.Background(), "user1", "hey!", domain.Marketing))

		rateLimitHandler.AssertNotCalled(t, "IsRateLimited")
		mailer.AssertNotCalled(t, "SendEmail")
	})
}
