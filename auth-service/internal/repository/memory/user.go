package memory

import (
	"authjwt/internal/domain"
	"context"
	"github.com/pkg/errors"
	"sync"
)

type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepo) CreateUser(_ context.Context, user *domain.User) error {
	//TODO::Я не стал делать по DDD - работаю с domain
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}

	//TODO::Я думаю тут подойдет и такие ошибки
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepo) ExistUser(ctx context.Context, userID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ID == userID {
			return true, nil
		}
	}

	return false, nil
}
