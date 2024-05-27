package inmemory

import (
	"context"
	"github.com/DedAzaMarks/ABS/internal/domain"
	myerrors "github.com/DedAzaMarks/ABS/internal/domain/errors"
	"sync"
)

type Cache struct {
	mu sync.RWMutex
	m  map[int64]*domain.User
}

func NewInmemory() *Cache {
	return &Cache{mu: sync.RWMutex{}, m: make(map[int64]*domain.User)}
}

func (c *Cache) GetUser(_ context.Context, userID int64) (*domain.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	user, ok := c.m[userID]
	if !ok {
		return nil, myerrors.ErrorUserNotFound
	}
	return user, nil
}

func (c *Cache) SetUser(_ context.Context, userID int64, user *domain.User) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[userID] = user
	return nil
}
