package repository

import (
	"fmt"
	"notification/internal/domain"
)

var (
	// ErrRuleAlreadyExists is the error when the rate limit rule already exists in the data store.
	ErrRuleAlreadyExists = fmt.Errorf("rule already exists")
)

// RateLimitRuleRepository is the abstract representation of the rate limit rule repository.
type RateLimitRuleRepository interface {
	// GetByNotificationType retrieves a rate limit rule by notification type.
	GetByNotificationType(notificationType domain.NotificationType) (domain.RateLimitRule, error)
	// Save stores a rate limit rule in the repository.
	Save(notificationType domain.NotificationType, rule domain.RateLimitRule) error
}

// NewInMemoryRateLimitRuleRepository creates a new InMemoryRateLimitRuleRepository instance.
func NewInMemoryRateLimitRuleRepository() *InMemoryRateLimitRuleRepository {
	return &InMemoryRateLimitRuleRepository{
		rules: make(domain.RateLimitRules),
	}
}

// InMemoryRateLimitRuleRepository is the in-memory representation of the rate limit rule repository.
type InMemoryRateLimitRuleRepository struct {
	rules domain.RateLimitRules
}

// GetByNotificationType retrieves a rate limit rule by notification type.
func (i InMemoryRateLimitRuleRepository) GetByNotificationType(
	notificationType domain.NotificationType) (domain.RateLimitRule, error) {
	return i.rules[notificationType], nil
}

// Save stores a rate limit rule in the repository.
func (i InMemoryRateLimitRuleRepository) Save(notificationType domain.NotificationType,
	rule domain.RateLimitRule) error {
	if _, ok := i.rules[notificationType]; ok {
		return ErrRuleAlreadyExists
	}

	i.rules[notificationType] = rule

	return nil
}
