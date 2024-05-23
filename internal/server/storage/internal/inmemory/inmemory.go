package inmemory

import (
	"context"
	"errors"
	"fmt"
	"github.com/DedAzaMarks/ABS/internal/domain"
	myerrors "github.com/DedAzaMarks/ABS/internal/domain/errors"
	"github.com/google/uuid"
	"log"
	"sync"
)

type Inmemory struct {
	mu    sync.RWMutex
	users map[string]*domain.TGUser
}

func NewInmemory() (*Inmemory, error) {
	return &Inmemory{users: make(map[string]*domain.TGUser)}, nil
}

func (s *Inmemory) GetUser(ctx context.Context, userID string) (*domain.TGUser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *Inmemory) AddNewUser(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userID]; ok {
		log.Print("user already exists")
		return fmt.Errorf("%w: %s", myerrors.ErrorUserAlreadyExists, userID)
	}
	s.users[userID] = domain.NewTGUser(userID)
	log.Print("user created")
	return nil
}

func (s *Inmemory) AddNewClient(ctx context.Context, userID, clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userID]; !ok {
		log.Print("user not found")
		return errors.New("user not found")
	}
	s.users[userID].Client = uuid.MustParse(clientID)
	log.Print("new client added:", clientID)
	return nil
}
