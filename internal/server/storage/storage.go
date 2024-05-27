package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/DedAzaMarks/ABS/internal/domain"
	myerrors "github.com/DedAzaMarks/ABS/internal/domain/errors"
	"github.com/DedAzaMarks/ABS/internal/server/storage/cache"
	"github.com/DedAzaMarks/ABS/internal/server/storage/persistent"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Storage interface {
	SaveUser(ctx context.Context, userID int64, user *domain.User) error
	LoadUser(ctx context.Context, userID int64) (*domain.User, error)
	GetUsersByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]*domain.User, error)

	AddNewDevice(ctx context.Context, userID int64, deviceID uuid.UUID, deviceName string) error
}

type storage struct {
	db persistent.Repository
	c  cache.Cache
}

func NewCachedStorage(ctx context.Context, repoType persistent.Type, cacheType cache.Type, client *redis.Client) (Storage, error) {
	db, err := persistent.GetRepo(ctx, repoType)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}
	c, err := cache.GetCache(ctx, cacheType, client)
	return &storage{
		db: db,
		c:  c,
	}, nil
}

func (s *storage) SaveUser(ctx context.Context, userID int64, user *domain.User) error {
	if err := s.c.SetUser(ctx, userID, user); err != nil {
		if !errors.Is(err, myerrors.ErrorUserNotFound) {
			return fmt.Errorf("failed to set user: %w", err)
		}
		return fmt.Errorf("failed to save user: %w", err)
	}
	dto := domain.TGUser2DTO(user)
	if err := s.db.SaveUser(ctx, dto); err != nil {
		if !errors.Is(err, myerrors.ErrorUserAlreadyExists) {
			return fmt.Errorf("failed to set user in persistent storage: %w", err)
		}
	}
	return nil
}

func (s *storage) LoadUser(ctx context.Context, userID int64) (*domain.User, error) {
	u, err := s.c.GetUser(ctx, userID)
	if err == nil {
		return u, nil
	}
	if !errors.Is(err, myerrors.ErrorUserNotFound) {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	userDTO, err := s.db.LoadUser(ctx, userID)
	if err != nil {
		if errors.Is(err, myerrors.ErrorUserNotFound) {
			return nil, fmt.Errorf("user is apsent in persistnt storage %w", err)
		} else if errors.Is(err, myerrors.ErrorDeviceNotFound) {
		} else {
			return nil, fmt.Errorf("failed to load user: %w", err)
		}
	}
	res := domain.DTO2TGUser(userDTO)
	if err := s.c.SetUser(ctx, userID, res); err != nil {
		return nil, fmt.Errorf("failed to set user in cache: %w", err)
	}
	return res, nil
}

func (s *storage) GetUsersByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]*domain.User, error) {
	us, err := s.db.GetUsersByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by device ID: %w", err)
	}
	if us != nil {
		return func() []*domain.User {
			res := make([]*domain.User, 0, len(us))
			for _, uDTO := range us {
				res = append(res, domain.DTO2TGUser(uDTO))
			}
			return res
		}(), nil
	}
	dtos, err := s.db.GetUsersByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by device ID: %w", err)
	}
	res := make([]*domain.User, 0, len(dtos))
	for _, dto := range dtos {
		res = append(res, domain.DTO2TGUser(dto))
	}
	return res, nil
}

func (s *storage) AddNewDevice(ctx context.Context, userID int64, deviceID uuid.UUID, deviceName string) error {
	u, err := s.c.GetUser(ctx, userID)
	if err != nil {
		if !errors.Is(err, myerrors.ErrorUserNotFound) {
			return fmt.Errorf("failed to get user: %w", err)
		}
	}
	if u == nil {
		dto, err := s.db.LoadUser(ctx, userID)
		if err != nil {
			if errors.Is(err, myerrors.ErrorUserNotFound) {
				return fmt.Errorf("user is apsent in persistnt storage %w", err)
			} else if errors.Is(err, myerrors.ErrorDeviceNotFound) {
			} else {
				return fmt.Errorf("failed to load user: %w", err)
			}
		}
		u = domain.DTO2TGUser(dto)
	}
	u.Devices = append(u.Devices, domain.Device{
		ID:   deviceID,
		Name: deviceName,
	})
	if err := s.c.SetUser(ctx, userID, u); err != nil {
		return fmt.Errorf("failed to set user in cache: %w", err)
	}
	if err := s.db.AddNewDevice(ctx, userID, deviceID, deviceName); err != nil {
		return fmt.Errorf("failed to add device in persistent storage: %w", err)
	}
	dto, err := s.db.LoadUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to load user: %w", err)
	}
	user := domain.DTO2TGUser(dto)
	if err := s.c.SetUser(ctx, userID, user); err != nil {
		return err
	}
	return nil
}
