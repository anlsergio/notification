package service

import (
	"context"
	"fmt"
	"notification/internal/domain"
	"notification/internal/repository"
)

// NotificationSender is the abstract representation of the NotificationSender service layer.
type NotificationSender interface {
	// Send sends a message to the given user depending on the notification type.
	Send(ctx context.Context, userID string, msg string, notificationType domain.NotificationType) error
}

// NewEmailNotificationSender creates a new EmailNotificationSender instance.
func NewEmailNotificationSender(rateLimitHandler RateLimitHandler,
	mailClient Mailer,
	userRepo repository.UserRepository) *EmailNotificationSender {
	return &EmailNotificationSender{
		rateLimitHandler: rateLimitHandler,
		client:           mailClient,
		userRepo:         userRepo,
	}
}

// EmailNotificationSender is the concrete email notification sender.
type EmailNotificationSender struct {
	rateLimitHandler RateLimitHandler
	client           Mailer
	userRepo         repository.UserRepository
}

// Send sends an email notification message to the given user depending on the notification type.
func (e EmailNotificationSender) Send(ctx context.Context,
	userID string, msg string, notificationType domain.NotificationType) error {
	ok, err := e.rateLimitHandler.Check(ctx, userID, notificationType)
	if err != nil {
		return fmt.Errorf("rate limit check fail: %w", err)
	}
	if !ok {
		// TODO: custom error type for proper error assertion.
		return fmt.Errorf("notification type %s exceeds the rate limit", notificationType)
	}

	user, err := e.userRepo.Get(userID)
	if err != nil {
		return fmt.Errorf("get user fail: %w", err)
	}

	subject := e.defineSubject(notificationType)
	if err := e.client.SendEmail(user.Email, subject, msg); err != nil {
		return fmt.Errorf("send mail fail: %w", err)
	}

	// TODO: retry in error?
	return e.rateLimitHandler.IncrementCount(ctx, userID, notificationType)
}

func (e EmailNotificationSender) defineSubject(notificationType domain.NotificationType) string {
	switch notificationType {
	case domain.Status:
		return fmt.Sprintf("%s: there's a new status update", notificationType)
	case domain.Marketing:
		return fmt.Sprintf("%s: we've got a new offer for you!", notificationType)
	case domain.News:
		return fmt.Sprintf("%s: we've got some news for you!", notificationType)
	default:
		return fmt.Sprintf("Notification")
	}
}
