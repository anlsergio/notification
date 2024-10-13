package service_test

import (
	"context"
	"errors"
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
			On("LockIfAvailable", mock.Anything, mock.Anything, mock.Anything).
			Return(time.Duration(0), nil, nil)

		mailer := mocks.NewMailSender(t)
		mailer.
			On("SendEmail", mock.Anything, mock.Anything).
			Return(nil)

		userRepo := mocks.NewUserRepository(t)
		userRepo.
			On("Get", mock.Anything).
			Return(domain.User{}, nil)

		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("")
		cacheSvc.
			On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		notification := domain.Notification{
			CorrelationID: "0990cc56-f1b7-4f69-bc60-08fac22d41bd",
			Type:          domain.Marketing,
			Message:       "Hey there!",
		}

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo, cacheSvc)
		_, err := svc.Send(context.Background(), "user1", notification)
		assert.NoError(t, err)
	})

	t.Run("rate limit exceeded", func(t *testing.T) {
		retryAfter := time.Minute

		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("LockIfAvailable", mock.Anything, mock.Anything, mock.Anything).
			Return(retryAfter, nil, service.ErrRateLimitExceeded).
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

		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("")

		notification := domain.Notification{
			CorrelationID: "0990cc56-f1b7-4f69-bc60-08fac22d41bd",
			Type:          domain.Marketing,
			Message:       "Hey there!",
		}

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo, cacheSvc)
		gotRetryAfter, err := svc.Send(context.Background(), "user1", notification)
		assert.ErrorIs(t, err, service.ErrRateLimitExceeded)
		assert.Equal(t, retryAfter, gotRetryAfter)

		mailer.AssertNotCalled(t, "SendEmail")
		cacheSvc.AssertNotCalled(t, "Set")
	})

	t.Run("invalid user", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("LockIfAvailable", mock.Anything, mock.Anything, mock.Anything).
			Return(time.Duration(0), nil, nil).
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

		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("")

		notification := domain.Notification{
			CorrelationID: "0990cc56-f1b7-4f69-bc60-08fac22d41bd",
			Type:          domain.Marketing,
			Message:       "Hey there!",
		}

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo, cacheSvc)
		_, err := svc.Send(context.Background(), "user1", notification)
		assert.Error(t, err)

		rateLimitHandler.AssertNotCalled(t, "LockIfAvailable")
		mailer.AssertNotCalled(t, "SendEmail")
		cacheSvc.AssertNotCalled(t, "Set")
	})

	t.Run("notification has already been processed", func(t *testing.T) {
		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("LockIfAvailable", mock.Anything, mock.Anything, mock.Anything).
			Return(time.Duration(0), nil, nil).
			Maybe()

		mailer := mocks.NewMailSender(t)
		mailer.
			On("SendEmail", mock.Anything, mock.Anything).
			Return(nil).
			Maybe()

		userRepo := mocks.NewUserRepository(t)
		userRepo.
			On("Get", mock.Anything).
			Return(domain.User{}, nil).
			Maybe()

		// if a corresponding key exists in cache, it means the same notification
		// has already been sent.
		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("processed")

		notification := domain.Notification{
			CorrelationID: "0990cc56-f1b7-4f69-bc60-08fac22d41bd",
			Type:          domain.Marketing,
			Message:       "Hey there!",
		}

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo, cacheSvc)
		_, err := svc.Send(context.Background(), "user1", notification)
		assert.ErrorIs(t, err, service.ErrIdempotencyViolation)

		rateLimitHandler.AssertNotCalled(t, "LockIfAvailable")
		mailer.AssertNotCalled(t, "SendEmail")
		cacheSvc.AssertNotCalled(t, "Set")
	})

	t.Run("release rate-limiting lock", func(t *testing.T) {
		var rolledBack bool

		spyRollback := func() error {
			rolledBack = true
			return nil
		}

		rateLimitHandler := mocks.NewRateLimitHandler(t)
		rateLimitHandler.
			On("LockIfAvailable", mock.Anything, mock.Anything, mock.Anything).
			Return(time.Duration(0), spyRollback, nil)

		mailer := mocks.NewMailSender(t)
		mailer.
			On("SendEmail", mock.Anything, mock.Anything).
			Return(errors.New("oops"))

		userRepo := mocks.NewUserRepository(t)
		userRepo.
			On("Get", mock.Anything).
			Return(domain.User{}, nil)

		cacheSvc := mocks.NewCache(t)
		cacheSvc.
			On("Get", mock.Anything, mock.Anything).
			Return("")
		cacheSvc.
			On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Maybe()

		notification := domain.Notification{
			CorrelationID: "0990cc56-f1b7-4f69-bc60-08fac22d41bd",
			Type:          domain.Marketing,
			Message:       "Hey there!",
		}

		svc := service.NewEmailNotificationSender(rateLimitHandler, mailer, userRepo, cacheSvc)
		_, err := svc.Send(context.Background(), "user1", notification)
		assert.Error(t, err)

		t.Run("lock is rolled back", func(t *testing.T) {
			assert.True(t, rolledBack)
		})

		t.Run("notification is not marked as processed", func(t *testing.T) {
			cacheSvc.AssertNotCalled(t, "Set")
		})
	})
}
