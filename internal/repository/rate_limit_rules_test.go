package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"notification/internal/domain"
	"notification/internal/repository"
	"testing"
	"time"
)

func TestInMemoryRateLimitRuleRepository_GetByNotificationType(t *testing.T) {
	rule := domain.RateLimitRule{
		MaxCount:   5,
		Expiration: time.Minute,
	}

	repo := repository.NewInMemoryRateLimitRuleRepository()
	require.NoError(t, repo.Save(domain.Marketing, rule))
	got, err := repo.GetByNotificationType(domain.Marketing)
	require.NoError(t, err)
	require.Equal(t, rule, got)
}

func TestInMemoryRateLimitRuleRepository_Save(t *testing.T) {
	t.Run("saves successfully", func(t *testing.T) {
		rule := domain.RateLimitRule{
			MaxCount:   5,
			Expiration: time.Minute,
		}

		repo := repository.NewInMemoryRateLimitRuleRepository()
		require.NoError(t, repo.Save(domain.Marketing, rule))
		got, err := repo.GetByNotificationType(domain.Marketing)
		require.NoError(t, err)
		require.Equal(t, rule, got)
	})

	t.Run("duplicate rules for the same notification type not allowed", func(t *testing.T) {
		rule1 := domain.RateLimitRule{
			MaxCount:   5,
			Expiration: time.Minute,
		}

		repo := repository.NewInMemoryRateLimitRuleRepository()
		require.NoError(t, repo.Save(domain.Marketing, rule1))

		// tries saving a new rule for the same notification type
		err := repo.Save(domain.Marketing, domain.RateLimitRule{})
		assert.ErrorIs(t, repository.ErrRuleAlreadyExists, err)

		// original rule configuration is not affected by the last saving attempt.
		got, err := repo.GetByNotificationType(domain.Marketing)
		require.NoError(t, err)
		require.Equal(t, rule1, got)
	})
}
