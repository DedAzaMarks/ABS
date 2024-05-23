package storage

import (
	"context"
	"github.com/DedAzaMarks/ABS/internal/server/storage/internal/inmemory"
	"github.com/DedAzaMarks/ABS/internal/server/storage/internal/postgres"
)

func NewInmemory() (Repo, error) {
	return inmemory.NewInmemory()
}

func NewPostgres(ctx context.Context, dsn string) (Repo, error) {
	return postgres.NewPostgres(ctx, dsn)
}
