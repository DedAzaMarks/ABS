package storage

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/DedAzaMarks/ABS/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrorUserAlreadyExists = errors.New("user with this id already exists")
)

type _ interface {
	GetById(uuid.UUID)
	GetAll()
	Insert()
	Delete()
}

type Repo interface {
	AddNewUser(userID string) error
	AddNewClient(userID, clientID string) error
	GetUserClients(userID string) ([]uuid.UUID, error)
}

func GetRepo(repoType string) (Repo, error) {
	switch repoType {
	case "inmemory":
		return &inmemory{users: make(map[string]*domain.TGUser)}, nil
	case "postgres":
		return nil, errors.New("not implemented")
	default:
		return nil, errors.New("repo type Not supported")
	}
}

type inmemory struct {
	mu    sync.RWMutex
	users map[string]*domain.TGUser
}

func (s *inmemory) GetUserClients(userID string) ([]uuid.UUID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user.Clients, nil
}

func (s *inmemory) AddNewUser(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userID]; ok {
		log.Print("user already exists")
		return fmt.Errorf("%w: %s", ErrorUserAlreadyExists, userID)
	}
	s.users[userID] = &domain.TGUser{
		State:   domain.EMPTY,
		UserID:  userID,
		Clients: []uuid.UUID{},
	}
	log.Print("user created")
	return nil
}

func (s *inmemory) AddNewClient(userID, clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userID]; !ok {
		log.Print("user not found")
		return errors.New("user not found")
	}
	s.users[userID].Clients = append(s.users[userID].Clients, uuid.MustParse(clientID))
	log.Print("new client added:", clientID)
	return nil
}
