package cache

import (
	"context"
	"github.com/DedAzaMarks/ABS/internal/server/storage/cache/internal/inmemory"
	rds "github.com/DedAzaMarks/ABS/internal/server/storage/cache/internal/redis"

	"github.com/go-redis/redis/v8"
)

func NewInmemory() (Cache, error) {
	return inmemory.NewInmemory(), nil
}

func NewRedis(ctx context.Context, redis *redis.Client) (Cache, error) {
	return rds.New(ctx, redis), nil
}
