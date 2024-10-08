package service

// NotificationType defines the different notification types.
type NotificationType int

const (
	// Status represents the notification for status updates.
	Status NotificationType = iota + 1
	// News represents the notification for news about the product.
	News
	// Marketing represents our marketing campaign notifications.
	Marketing
)

// Notification is the abstract representation of the Notification service layer.
type Notification interface {
	// Send sends a message to the given user depending on the notification type.
	Send(userID string, msg string, notificationType NotificationType) error
}

func NewEmailNotification(cacheService Cache) *EmailNotification {
	return &EmailNotification{
		cacheService: cacheService,
	}
}

type EmailNotification struct {
	cacheService Cache
}

func (e EmailNotification) Send(userID string, msg string, notificationType NotificationType) error {
	return nil
}
