package inmemory

import (
	"context"
	"errors"
	"fmt"
	"github.com/DedAzaMarks/ABS/internal/domain"
	myerrors "github.com/DedAzaMarks/ABS/internal/domain/errors"
	"github.com/google/uuid"
	"sync"
)

type user struct {
	userID     int64
	sessionKey string
	state      string
}

type userDeviceKey struct {
	userID   int64
	deviceID uuid.UUID
}

type device struct {
	deviceID   uuid.UUID
	deviceName string
}

type Inmemory struct {
	mu          sync.RWMutex
	users       map[int64]user
	userDevices []userDeviceKey
	devices     map[uuid.UUID]device
}

func NewInmemory() (*Inmemory, error) {
	return &Inmemory{
		users:   make(map[int64]user),
		devices: make(map[uuid.UUID]device),
	}, nil
}

func (s *Inmemory) SaveUser(_ context.Context, userDTO *domain.UserDTO) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userDTO.ID]; ok {
		return myerrors.ErrorUserAlreadyExists
	}
	s.users[userDTO.ID] = user{
		userID:     userDTO.ID,
		sessionKey: userDTO.SessionKey,
		state:      userDTO.State,
	}
	for _, d := range userDTO.Devices {
		s.userDevices = append(s.userDevices, userDeviceKey{
			userID:   userDTO.ID,
			deviceID: d.ID,
		})
	}
	for _, d := range userDTO.Devices {
		s.devices[d.ID] = device{
			deviceID:   d.ID,
			deviceName: d.Name,
		}
	}
	return nil
}

func (s *Inmemory) LoadUser(_ context.Context, userID int64) (*domain.UserDTO, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[userID]
	if !ok {
		return nil, myerrors.ErrorUserNotFound
	}
	res := &domain.UserDTO{
		ID:         u.userID,
		SessionKey: u.sessionKey,
		State:      u.state,
	}
	for _, ud := range s.userDevices {
		if ud.userID == u.userID {
			d, ok := s.devices[ud.deviceID]
			if !ok {
				return nil, fmt.Errorf("unregistered device: %s", ud.deviceID)
			}
			res.Devices = append(res.Devices, domain.DeviceDTO{
				ID:   d.deviceID,
				Name: d.deviceName,
			})
		}
	}
	return res, nil
}

func (s *Inmemory) GetUsersByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]*domain.UserDTO, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var m map[int64]*domain.UserDTO
	for _, userDevice := range s.userDevices {
		if userDevice.deviceID == deviceID {
			if _, ok := m[userDevice.userID]; ok {
				continue
			}
			userDTO, err := s.LoadUser(ctx, userDevice.userID)
			if err != nil {
				if errors.Is(err, myerrors.ErrorUserNotFound) {
					return nil, fmt.Errorf("device belongs to unknown user: %d", userDevice.userID)
				}
				return nil, fmt.Errorf("loading user: %w", err)
			}
			m[userDevice.userID] = userDTO
		}
	}
	var res []*domain.UserDTO
	for _, user := range m {
		res = append(res, user)
	}
	return res, nil
}

func (s *Inmemory) AddNewDevice(_ context.Context, userID int64, deviceID uuid.UUID, deviceName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.users[userID]
	if !ok {
		return myerrors.ErrorUserNotFound
	}
	s.devices[deviceID] = device{
		deviceID:   deviceID,
		deviceName: deviceName}
	s.userDevices = append(s.userDevices, userDeviceKey{
		userID:   userID,
		deviceID: deviceID})
	return nil
}
