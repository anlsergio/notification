package domain

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
