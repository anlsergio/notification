package service

import (
	"context"
	"fmt"
	"notification/internal/domain"
	"strconv"
	"time"
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
	key := fmt.Sprintf("%s:%s", userID, notificationType)
	rule := h.limitRules[notificationType]

	ok, err := h.checkAvailability(ctx, key, rule.MaxCount)
	if err != nil {
		return false, fmt.Errorf("check availability fail: %w", err)
	}

	if !ok {
		return false, nil
	}

	if err = h.incrementCount(ctx, key, rule.Expiration); err != nil {
		return false, fmt.Errorf("increment count fail: %w", err)
	}

	return true, nil
}

// check returns True if there's capacity available for the notification
// to be sent based on the maximum allowed count for the given key.
func (h CacheRateLimitHandler) checkAvailability(ctx context.Context,
	key string, maxCount int) (bool, error) {

	stringCounts := h.cacheService.Get(ctx, key)
	counts, err := strconv.Atoi(stringCounts)
	if err != nil {
		// if the int conversion fails and stringCounts is populated with anything but an empty string
		// at this point it's not safe to assume its correct int counterpart.
		if stringCounts != "" {
			return false, fmt.Errorf("failed converting notification counts from cache: %w", err)
		}
		// if it's an empty string, it's probably because the key doesn't currently exist in the cache yet
		// and therefore, we can assume a 0 count.
		counts = 0
	}

	if counts >= maxCount {
		return false, nil
	}

	return true, nil
}

// incrementCount adds to the rate limit counter based the key, applying the specified TTL.
func (h CacheRateLimitHandler) incrementCount(ctx context.Context,
	key string, ttl time.Duration) error {
	return h.cacheService.Incr(ctx, key, ttl)
}
