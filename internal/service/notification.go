package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"notification/internal/domain"
	"notification/internal/repository"
	"time"
)

var (
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
	// idempotency check: ensures that the notification hasn't already been processed.
	if e.isAlreadyProcessed(ctx, notification.CorrelationID) {
		return 0, newIdempotencyError(notification.CorrelationID)
	}

	user, err := e.userRepo.Get(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user: %w", err)
	}

	lockResult, err := e.acquireRateLimitLock(ctx, userID, notification.Type)
	if err != nil {
		if lockResult != nil {
			retryAfter = lockResult.RetryAfter
		}
		return retryAfter, err
	}

	subject := e.defineSubject(notification.Type)
	if err := e.client.SendEmail(user.Email, subject, notification.Message); err != nil {
		// if the email could not be sent for any reason, release the rate-limit lock.
		e.safeRollback(lockResult)
		return 0, fmt.Errorf("failed to send email: %w", err)
	}

	// If everything went fine, mark the current notification as processed
	// for the idempotency check.
	if err = e.markAsProcessed(ctx, notification.CorrelationID, time.Hour*24); err != nil {
		return 0, fmt.Errorf("failed to mark notification as processed: %w", err)
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

func newIdempotencyError(correlationID string) error {
	return errors.Join(ErrIdempotencyViolation,
		fmt.Errorf("the notification of correlation ID %s has already been processed", correlationID))
}

func (e EmailNotificationSender) acquireRateLimitLock(ctx context.Context,
	userID string, notificationType domain.NotificationType) (*LockResult, error) {

	lockResult, err := e.rateLimitHandler.LockIfAvailable(ctx, userID, notificationType)
	if err != nil {
		if errors.Is(err, ErrRateLimitExceeded) {
			return lockResult, fmt.Errorf("notification type %s exceeds the rate limit: %w", notificationType, err)
		}
		return nil, fmt.Errorf("rate limit check fail: %w", err)
	}
	return lockResult, nil
}

func (e EmailNotificationSender) safeRollback(lockResult *LockResult) {
	if lockResult == nil {
		// no need to roll back because the lock hasn't been acquired.
		return
	}

	if err := lockResult.Rollback(); err != nil {
		log.Printf("rollback of rate-limit counter failed: %v", err)
	}
}
