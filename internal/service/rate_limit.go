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
	// IsRateLimited returns True if there's capacity available for the notification
	// to be sent for the given user.
	IsRateLimited(ctx context.Context, userID string, notificationType domain.NotificationType) (bool, error)
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

func (h CacheRateLimitHandler) IsRateLimited(ctx context.Context,
	userID string, notificationType domain.NotificationType) (bool, error) {

	ok, err := h.checkAvailability(ctx, userID, notificationType)
	if err != nil {
		return false, fmt.Errorf("check availability fail: %w", err)
	}

	if err = h.incrementCount(ctx, userID, notificationType); err != nil {
		return false, fmt.Errorf("increment count fail: %w", err)
	}

	return ok, nil
}

// check returns True if there's capacity available for the notification
// to be sent for the given user.
// TODO: Handle concurrent checks where depending on the number of replicas
// the notification system might misbehave, allowing more notifications than it should.
// Transactional guarantee? Locking someway?
func (h CacheRateLimitHandler) checkAvailability(ctx context.Context,
	userID string, notificationType domain.NotificationType) (bool, error) {
	cacheKey := fmt.Sprintf("%s:%s", userID, notificationType)

	stringCounts := h.cacheService.Get(ctx, cacheKey)
	counts, err := strconv.Atoi(stringCounts)
	if err != nil {
		if stringCounts != "" {
			return false, fmt.Errorf("failed converting notification counts from cache: %w", err)
		}
		counts = 0
	}

	rule := h.limitRules[notificationType]
	if counts >= rule.MaxCount {
		return false, nil
	}

	return true, nil
}

// incrementCount adds to the rate limit counter for the given userID + notification type combination.
func (h CacheRateLimitHandler) incrementCount(ctx context.Context,
	userID string, notificationType domain.NotificationType) error {
	cacheKey := fmt.Sprintf("%s:%s", userID, notificationType)
	rule := h.limitRules[notificationType]
	return h.cacheService.Incr(ctx, cacheKey, rule.Expiration)
}
