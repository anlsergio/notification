package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"notification/internal/domain"
	"notification/internal/repository"
	"time"
)

var (
	// ErrRateLimitExceeded is the error when the notification cannot be sent because
	// it exceeds the rate limiting rules defined.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	// ErrIdempotencyViolation is the error when the same notification has already been processed before.
	ErrIdempotencyViolation = errors.New("notification already processed")
)

// NotificationSender is the abstract representation of the NotificationSender service layer.
type NotificationSender interface {
	// Send sends a message to the given user depending on the notification type.
	Send(ctx context.Context,
		userID string, notification domain.Notification) (retryAfter time.Duration, err error)
}

// NewEmailNotificationSender creates a new EmailNotificationSender instance.
func NewEmailNotificationSender(rateLimitHandler RateLimitHandler,
	mailClient Mailer,
	userRepo repository.UserRepository,
	cacheService Cache) *EmailNotificationSender {
	return &EmailNotificationSender{
		rateLimitHandler: rateLimitHandler,
		client:           mailClient,
		userRepo:         userRepo,
		cache:            cacheService,
	}
}

// EmailNotificationSender is the concrete email notification sender.
type EmailNotificationSender struct {
	rateLimitHandler RateLimitHandler
	client           Mailer
	userRepo         repository.UserRepository
	cache            Cache
}

// Send sends an email notification message to the given user depending on the notification type.
// It returns ErrRateLimitExceeded if the notification being sent exceeds the pre-defined rate-limiting rules.
func (e EmailNotificationSender) Send(ctx context.Context,
	userID string, notification domain.Notification) (retryAfter time.Duration, err error) {
	// idempotency check
	if e.isAlreadyProcessed(ctx, notification.CorrelationID) {
		return 0, errors.Join(ErrIdempotencyViolation,
			fmt.Errorf("the notification of correlation ID %s has already been processed", notification.CorrelationID))
	}

	user, err := e.userRepo.Get(userID)
	if err != nil {
		return 0, fmt.Errorf("get user fail: %w", err)
	}

	ok, retryAfter, err := e.rateLimitHandler.IsRateLimited(ctx, userID, notification.Type)
	if err != nil {
		return 0, fmt.Errorf("rate limit check fail: %w", err)
	}
	if !ok {
		return retryAfter, errors.Join(ErrRateLimitExceeded,
			fmt.Errorf("notification type %s exceeds the rate limit", notification.Type))
	}

	subject := e.defineSubject(notification.Type)
	if err := e.client.SendEmail(user.Email, subject, notification.Message); err != nil {
		return 0, fmt.Errorf("send mail fail: %w", err)
	}

	// If everything went fine, mark the current notification as processed
	// for the idempotency check.
	if err = e.markAsProcessed(ctx, notification.CorrelationID, time.Hour*24); err != nil {
		return 0, fmt.Errorf("mark as processed fail: %w", err)
	}

	return 0, nil
}

func (e EmailNotificationSender) defineSubject(notificationType domain.NotificationType) string {
	var subject string
	switch notificationType {
	case domain.Status:
		subject = "there's a new status update"
	case domain.Marketing:
		subject = "we've got a new offer for you!"
	case domain.News:
		subject = "we've got some news for you!"
	default:
		return "Notification"
	}

	prefix := cases.
		Title(language.English, cases.Compact).
		String(notificationType.String())

	return fmt.Sprintf("%s: %s", prefix, subject)
}

func (e EmailNotificationSender) isAlreadyProcessed(ctx context.Context, correlationID string) bool {
	return e.cache.Get(ctx, correlationID) != ""
}

func (e EmailNotificationSender) markAsProcessed(ctx context.Context,
	correlationID string, expiration time.Duration) error {
	return e.cache.Set(ctx, correlationID, "processed", expiration)
}
