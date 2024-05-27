package persistent

import (
	"context"
	"github.com/DedAzaMarks/ABS/internal/server/storage/persistent/internal/inmemory"
	"github.com/DedAzaMarks/ABS/internal/server/storage/persistent/internal/postgres"
)

func NewInmemory() (Repository, error) {
	return inmemory.NewInmemory()
}

func NewPostgres(ctx context.Context) (Repository, error) {
	return postgres.NewPostgres(ctx)
}
