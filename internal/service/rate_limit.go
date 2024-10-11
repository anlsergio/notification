package service

import (
	"context"
	"fmt"
	"notification/internal/domain"
	"strconv"
)

// RateLimitHandler is the abstract representation of the rate limit checker,
// responsible for informing if there's capacity available for the notification to be sent
// to a given user using a Leaky Bucket algorithm.
type RateLimitHandler interface {
	// Check returns True if there's capacity available for the notification
	// to be sent for the given user.
	Check(ctx context.Context, userID string, notificationType domain.NotificationType) (bool, error)
	// IncrementCount TODO: add go doc.
	IncrementCount(ctx context.Context, userID string, notificationType domain.NotificationType) error
}

// NewCacheRateLimitHandler creates a new CacheRateLimitHandler instance.
func NewCacheRateLimitHandler(cacheService Cache, rules domain.RateLimitRules) *CacheRateLimitHandler {
	return &CacheRateLimitHandler{
		cacheService: cacheService,
		limitRules:   rules,
	}
}

// CacheRateLimitHandler handles the rate limiting checks and state
// based on a cache service.
type CacheRateLimitHandler struct {
	cacheService Cache
	limitRules   domain.RateLimitRules
}

// Check returns True if there's capacity available for the notification
// to be sent for the given user.
// TODO: Handle concurrent checks where depending on the number of replicas
// the notification system might misbehave, allowing more notifications than it should.
// Transactional guarantee? Locking someway?
func (d CacheRateLimitHandler) Check(ctx context.Context,
	userID string, notificationType domain.NotificationType) (bool, error) {
	cacheKey := fmt.Sprintf("%s:%s", userID, notificationType)

	stringCounts := d.cacheService.Get(ctx, cacheKey)
	counts, err := strconv.Atoi(stringCounts)
	if err != nil {
		if stringCounts != "" {
			return false, fmt.Errorf("failed converting notification counts from cache: %w", err)
		}
		counts = 0
	}

	rule := d.limitRules[notificationType]
	if counts >= rule.MaxCount {
		return false, nil
	}

	return true, nil
}

// IncrementCount adds to the rate limit counter for the given userID + notification type combination.
func (d CacheRateLimitHandler) IncrementCount(ctx context.Context,
	userID string, notificationType domain.NotificationType) error {
	cacheKey := fmt.Sprintf("%s:%s", userID, notificationType)
	rule := d.limitRules[notificationType]
	return d.cacheService.Incr(ctx, cacheKey, rule.Expiration)
}
