package domain

import "time"

// RateLimitRules defines the rate limit rules for a given notification type.
type RateLimitRules map[NotificationType]RateLimitRule

// RateLimitRule defines the rate limit rule configuration.
type RateLimitRule struct {
	// MaxCount is the max notification count allowed for a given time span.
	MaxCount int
	// Expiration is the time span defined for limiting a certain number of messages.
	Expiration time.Duration
}
