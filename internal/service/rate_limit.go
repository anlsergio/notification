package service

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// RateLimitHandler is the abstract representation of the rate limit checker,
// responsible for informing if there's capacity available for the notification to be sent
// to a given user.
type RateLimitHandler interface {
	// Check returns True if there's capacity available for the notification
	// to be sent for the given user.
	Check(ctx context.Context, userID string, notificationType NotificationType) (bool, error)
	// IncrementCount TODO: add go doc.
	IncrementCount(ctx context.Context, userID string, notificationType NotificationType) error
}

// RateLimitRules defines the rate limit rules for a given notification type.
type RateLimitRules map[NotificationType]RateLimitRule

// RateLimitRule defines the rate limit rule configuration.
type RateLimitRule struct {
	// MaxCount is the max notification count allowed for a given time span.
	MaxCount int
	// Expiration is the time span defined for limiting a certain number of messages.
	Expiration time.Duration
}

// NewCacheRateLimitHandler creates a new CacheRateLimitHandler instance.
// TODO: the config rules should probably be injected into the constructor, instead
// of being hard-coded, to give more control to tests.
func NewCacheRateLimitHandler(cacheService Cache) *CacheRateLimitHandler {
	// TODO: this set of configurations could be fetched
	// from a config service.
	rules := RateLimitRules{
		Status: RateLimitRule{
			2,
			time.Minute * 1,
		},
		News: RateLimitRule{
			1,
			time.Hour * 24,
		},
		Marketing: RateLimitRule{
			3,
			time.Hour * 1,
		},
	}

	return &CacheRateLimitHandler{
		cacheService: cacheService,
		limitRules:   rules,
	}
}

// CacheRateLimitHandler handles the rate limiting checks and state
// based on a cache service.
type CacheRateLimitHandler struct {
	cacheService Cache
	limitRules   RateLimitRules
}

// Check returns True if there's capacity available for the notification
// to be sent for the given user.
// TODO: Handle concurrent checks where depending on the number of replicas
// the notification system might misbehave, allowing more notifications than it should.
// Transactional guarantee? Locking someway?
func (d CacheRateLimitHandler) Check(ctx context.Context, userID string, notificationType NotificationType) (bool, error) {
	cacheKey := fmt.Sprintf("%s:%s", userID, notificationType)

	notificationCounts := d.cacheService.Get(ctx, cacheKey)
	counts, err := strconv.Atoi(notificationCounts)
	if err != nil {
		return false, fmt.Errorf("failed converting notification counts from cache: %w", err)
	}

	rule := d.limitRules[notificationType]
	if counts >= rule.MaxCount {
		return false, nil
	}

	return true, nil
}

// IncrementCount adds to the rate limit counter for the given userID + notification type combination.
func (d CacheRateLimitHandler) IncrementCount(ctx context.Context, userID string, notificationType NotificationType) error {
	cacheKey := fmt.Sprintf("%s:%s", userID, notificationType)
	rule := d.limitRules[notificationType]
	return d.cacheService.Incr(ctx, cacheKey, rule.Expiration)
}
