package service

import (
	"context"
	"fmt"
	"time"
)

// RateLimitHandler is the abstract representation of the rate limit checker,
// responsible for informing if there's capacity available for the notification to be sent
// to a given user.
type RateLimitHandler interface {
	// Check returns True if there's capacity available for the notification
	// to be sent for the given user.
	Check(ctx context.Context, userID string, notificationType NotificationType) (bool, error)
	IncrementCount(ctx context.Context, userID string, notificationType NotificationType) error
}

// RateLimitConfig defines the rate limit configurations for notifications.
type RateLimitConfig struct {
	// NotificationType is the notification type this config has effect upon.
	NotificationType NotificationType
	// MaxCount is the max notification count allowed for a given time span.
	MaxCount int
	// TimeSpan is the time span defined for limiting a certain number of messages.
	TimeSpan time.Duration
}

func NewCacheRateLimitChecker(cacheService Cache) *CacheRateLimitHandler {
	// TODO: this set of configurations could be fetched
	// from a config service.
	configs := []RateLimitConfig{
		{
			Status,
			2,
			time.Minute * 1,
		},
		{
			News,
			1,
			time.Hour * 24,
		},
		{
			Marketing,
			3,
			time.Hour * 1,
		},
	}
	return &CacheRateLimitHandler{
		cacheService: cacheService,
		configs:      configs,
	}
}

type CacheRateLimitHandler struct {
	cacheService Cache
	configs      []RateLimitConfig
}

// Check returns True if there's capacity available for the notification
// to be sent for the given user.
// TODO: Handle concurrent checks where depending on the number of replicas
// the notification system might misbehave, allowing more notifications than it should.
// Transactional guarantee? Locking someway?
func (d CacheRateLimitHandler) Check(ctx context.Context, userID string, notificationType NotificationType) (bool, error) {
	cacheKey := fmt.Sprintf("%s:%s", userID, notificationType)

	notificationCounts, err := d.cacheService.Get(ctx, cacheKey)
	if err != nil {
		return false, fmt.Errorf("get notification count from cache: %w", err)
	}
	counts, ok := notificationCounts.(int)
	if !ok {
		return false, fmt.Errorf("failed to type cast notification count because it's not int")
	}

	// TODO: configs could be a hash table, identified by notificationType as key.
	for _, config := range d.configs {
		if notificationType == config.NotificationType {
			if counts >= config.MaxCount {
				return false, nil
			}
		}
	}

	return true, nil
}

// IncrementCount adds to the rate limit counter for the given userID + notification type combination.
func (d CacheRateLimitHandler) IncrementCount(ctx context.Context, userID string, notificationType NotificationType) error {
	cacheKey := fmt.Sprintf("%s:%s", userID, notificationType)
	// TODO: see how to increment by 1
	return d.cacheService.Set(ctx, cacheKey, 1)
}
