package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"notification/internal/domain"
	"notification/internal/repository"
	"testing"
)

func TestInMemoryUserRepository_Save(t *testing.T) {
	t.Run("user is saved", func(t *testing.T) {
		user := domain.User{
			ID:       "123-abc",
			Name:     "John",
			LastName: "Doe",
			Email:    "john.doe@example.com",
		}

		repo := repository.NewInMemoryUserRepository()
		require.NoError(t, repo.Save(user))

		savedUser, err := repo.Get(user.ID)
		require.NoError(t, err)
		require.NotEmpty(t, savedUser)
		require.Equal(t, user.ID, savedUser.ID)
	})

	t.Run("user is not saved: conflicting IDs", func(t *testing.T) {
		user := domain.User{
			ID:       "123-abc",
			Name:     "John",
			LastName: "Doe",
			Email:    "john.doe@example.com",
		}

		repo := repository.NewInMemoryUserRepository()
		require.NoError(t, repo.Save(user))

		err := repo.Save(domain.User{ID: "123-abc"})
		assert.ErrorIs(t, err, repository.ErrUserAlreadyExists)
	})

	t.Run("user is not saved: conflicting emails", func(t *testing.T) {
		user := domain.User{
			ID:       "123-abc",
			Name:     "John",
			LastName: "Doe",
			Email:    "john.doe@example.com",
		}

		repo := repository.NewInMemoryUserRepository()
		require.NoError(t, repo.Save(user))

		err := repo.Save(domain.User{Email: "john.doe@example.com"})
		assert.ErrorIs(t, err, repository.ErrUserAlreadyExists)
	})
}

func TestInMemoryUserRepository_Get(t *testing.T) {
	t.Run("user is found", func(t *testing.T) {
		user := domain.User{
			ID:       "123-abc",
			Name:     "John",
			LastName: "Doe",
			Email:    "john.doe@example.com",
		}

		repo := repository.NewInMemoryUserRepository()
		require.NoError(t, repo.Save(user))

		savedUser, err := repo.Get(user.ID)
		require.NoError(t, err)
		require.NotEmpty(t, savedUser)
		require.Equal(t, user.ID, savedUser.ID)
	})

	t.Run("user is not found", func(t *testing.T) {
		repo := repository.NewInMemoryUserRepository()
		_, err := repo.Get("invalid")
		assert.ErrorIs(t, repository.ErrInvalidUserID, err)
	})
}
