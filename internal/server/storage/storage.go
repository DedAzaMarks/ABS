package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/DedAzaMarks/ABS/internal/domain"
	"os"

	"github.com/google/uuid"
)

type _ interface {
	GetById(uuid.UUID)
	GetAll()
	Insert()
	Delete()
}

type Repo interface {
	AddNewUser(ctx context.Context, userID string) error
	AddNewClient(ctx context.Context, userID, clientID string) error
	GetUser(ctx context.Context, userID string) (*domain.TGUser, error)
}

type RepoType string

const (
	InMemory RepoType = "inmemory"
	Postgres RepoType = "postgres"
)

func GetRepo(ctx context.Context, repoType RepoType) (Repo, error) {
	switch repoType {
	case InMemory:
		return NewInmemory()
	case Postgres:
		postgresDSN := fmt.Sprintf(
			"host=%s port=%s database=%s user=%s password=%s",
			"postgres",
			"5432",
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
		)
		return NewPostgres(ctx, postgresDSN)
	default:
		return nil, errors.New("repo type Not supported")
	}
}
