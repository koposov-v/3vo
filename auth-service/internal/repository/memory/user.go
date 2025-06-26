package memory

import (
	"authjwt/internal/domain"
	"context"
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
