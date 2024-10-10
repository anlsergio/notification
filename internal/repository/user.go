package repository

import (
	"errors"
	"notification/internal/domain"
)

var (
	// ErrUserAlreadyExists is the error when a user is already present in the data store.
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserRepository is the abstract representation of the user repository.
// TODO: UserRepository is redundant with the package name.
type UserRepository interface {
	// Get retrieves a user by its ID.
	Get(id string) (*domain.User, error)
	// Save stores a given user in the repository.
	Save(user *domain.User) error
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository instance.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// InMemoryUserRepository is the in-memory representation of the user repository.
type InMemoryUserRepository struct {
	users map[string]*domain.User
}

// Get retrieves a user by its ID.
func (r *InMemoryUserRepository) Get(id string) (*domain.User, error) {
	return r.users[id], nil
}

// Save stores a given user in the repository.
func (r *InMemoryUserRepository) Save(user *domain.User) error {
	for _, u := range r.users {
		if u.Email == user.Email || u.ID == user.ID {
			return ErrUserAlreadyExists
		}
	}

	r.users[user.ID] = user

	return nil
}
