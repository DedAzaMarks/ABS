package persistent

import (
	"context"
	"errors"
	"github.com/DedAzaMarks/ABS/internal/domain"
	"github.com/google/uuid"
)

type Type string

const (
	InMemory Type = "inmemory"
	Postgres Type = "postgres"
)

type Repository interface {
	SaveUser(ctx context.Context, user *domain.UserDTO) error
	LoadUser(ctx context.Context, userID int64) (*domain.UserDTO, error)
	GetUsersByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]*domain.UserDTO, error)

	AddNewDevice(ctx context.Context, userID int64, deviceID uuid.UUID, deviceName string) error
}

func GetRepo(ctx context.Context, repoType Type) (Repository, error) {
	switch repoType {
	case InMemory:
		return NewInmemory()
	case Postgres:
		return NewPostgres(ctx)
	default:
		return nil, errors.New("repo type Not supported")
	}
}
