package service

import (
	"context"
	"fmt"
)

const (
	// Status represents the notification for status updates.
	Status NotificationType = iota + 1
	// News represents the notification for news about the product.
	News
	// Marketing represents our marketing campaign notifications.
	Marketing
)

// NotificationType defines the different notification types.
type NotificationType int

// String returns the string equivalent of NotificationType.
// It returns an empty string if the notification type is invalid.
func (t NotificationType) String() string {
	switch t {
	case Status:
		return "status"
	case News:
		return "news"
	case Marketing:
		return "marketing"
	default:
		return ""
	}
}

// NotificationSender is the abstract representation of the NotificationSender service layer.
type NotificationSender interface {
	// Send sends a message to the given user depending on the notification type.
	Send(ctx context.Context, userID string, msg string, notificationType NotificationType) error
}

// NewEmailNotificationSender creates a new EmailNotificationSender instance.
func NewEmailNotificationSender(rateLimitHandler RateLimitHandler, mailClient MailClient) *EmailNotificationSender {
	return &EmailNotificationSender{
		rateLimitHandler: rateLimitHandler,
		client:           mailClient,
	}
}

// EmailNotificationSender is the concrete email notification sender.
type EmailNotificationSender struct {
	rateLimitHandler RateLimitHandler
	client           MailClient
}

// Send sends an email notification message to the given user depending on the notification type.
func (e EmailNotificationSender) Send(ctx context.Context, userID string, msg string, notificationType NotificationType) error {
	ok, err := e.rateLimitHandler.Check(ctx, userID, notificationType)
	if err != nil {
		return fmt.Errorf("rate limit check fail: %w", err)
	}
	if !ok {
		// TODO: custom error type for proper error assertion.
		return fmt.Errorf("notification type %s exceeds the rate limit", notificationType)
	}

	if err := e.client.Send(); err != nil {
		return fmt.Errorf("send mail fail: %w", err)
	}

	// TODO: retry in error?
	return e.rateLimitHandler.IncrementCount(ctx, userID, notificationType)
}

// MailClient is the abstraction layer of the external email service integration itself.
type MailClient interface {
	// Send sends the email message through the appropriate external service integration.
	Send() error
}
