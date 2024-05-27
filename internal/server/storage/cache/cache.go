package cache

import (
	"context"
	"errors"
	"github.com/DedAzaMarks/ABS/internal/domain"
	"github.com/go-redis/redis/v8"
)

type Type string

const (
	InMemory Type = "inmemory"
	Redis    Type = "redis"
)

type Cache interface {
	GetUser(_ context.Context, userID int64) (*domain.User, error)
	SetUser(_ context.Context, userID int64, user *domain.User) error
}

func GetCache(ctx context.Context, cacheType Type, client *redis.Client) (Cache, error) {
	switch cacheType {
	case InMemory:
		return NewInmemory()
	case Redis:
		return NewRedis(ctx, client)
	default:
		return nil, errors.New("repo type Not supported")
	}
}
