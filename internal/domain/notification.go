package domain

import (
	"errors"
)

const (
	// Status represents the notification for status updates.
	Status NotificationType = iota + 1
	// News represents the notification for news about the product.
	News
	// Marketing represents our marketing campaign notifications.
	Marketing
)

var (
	// ErrInvalidNotificationType is the error when the provided notification type is invalid.
	ErrInvalidNotificationType = errors.New("unknown notification type")
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

// ToNotificationType converts a string into a corresponding NotificationType.
// It will error out if the string doesn't match any pre-defined notification type.
func ToNotificationType(s string) (NotificationType, error) {
	switch s {
	case "status":
		return Status, nil
	case "news":
		return News, nil
	case "marketing":
		return Marketing, nil
	default:
		return 0, ErrInvalidNotificationType
	}
}
