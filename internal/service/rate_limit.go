package service

import (
	"context"
	"errors"
	"fmt"
	"notification/internal/domain"
	"notification/internal/repository"
	"strconv"
	"time"
)

var (
	// ErrRateLimitExceeded is the error when the notification cannot be sent because
	// it exceeds the rate limiting rules defined.
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// RateLimitHandler is the abstract representation of the rate limit checker,
// responsible for informing if there's capacity available for the notification to be sent
// to a given user using a Leaky Bucket algorithm.
type RateLimitHandler interface {
	// IsRateLimited returns ErrRateLimitExceeded if there's no capacity available for the notification
	// to be sent for the given user.
	// If there's no capacity available it informs the caller through retryAfter
	// how much time is left until the next token is available.
	//
	// Important: for every time it returns stating that the capacity is available, it will
	// automatically increment the limit counter. Use the rollback function to handle the rollback
	// scenario, which is the caller's responsibility.
	IsRateLimited(ctx context.Context,
		userID string,
		notificationType domain.NotificationType) (retryAfter time.Duration, rollback func() error, err error)
}

// NewCacheRateLimitHandler creates a new CacheRateLimitHandler instance.
func NewCacheRateLimitHandler(cacheService Cache, rulesRepo repository.RateLimitRuleRepository) *CacheRateLimitHandler {
	return &CacheRateLimitHandler{
		cacheService: cacheService,
		repo:         rulesRepo,
	}
}

// CacheRateLimitHandler handles the rate limiting checks and state
// based on a cache service.
type CacheRateLimitHandler struct {
	cacheService Cache
	repo         repository.RateLimitRuleRepository
}

// IsRateLimited returns ErrRateLimitExceeded if there's no capacity available for the notification
// to be sent for the given user.
// If there's no capacity available it informs the caller through retryAfter
// how much time is left until the next token is available.
//
// Important: for every time it returns stating that the capacity is available, it will
// automatically increment the limit counter. Use the rollback function to handle the rollback
// scenario, which is the caller's responsibility.
func (h CacheRateLimitHandler) IsRateLimited(ctx context.Context,
	userID string, notificationType domain.NotificationType) (retryAfter time.Duration, rollback func() error, err error) {
	key := fmt.Sprintf("%s:%s", userID, notificationType)
	rule, err := h.repo.GetByNotificationType(notificationType)
	if err != nil {
		return 0, nil, fmt.Errorf("get rate limit rule by notification type fail: %w", err)
	}

	ok, err := h.checkAvailability(ctx, key, rule.MaxCount)
	if err != nil {
		return 0, nil, fmt.Errorf("check availability fail: %w", err)
	}

	if !ok {
		return rule.Expiration, nil, ErrRateLimitExceeded
	}

	if err = h.incrementCount(ctx, key, rule.Expiration); err != nil {
		return 0, nil, fmt.Errorf("increment count fail: %w", err)
	}

	// give the ability to roll back the operation to the caller.
	rollback = func() error {
		return h.cacheService.Decr(ctx, key)
	}

	return 0, rollback, nil
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
