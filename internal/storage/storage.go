package storage

import "github.com/google/uuid"

type Storage interface {
	GetById(uuid.UUID)
	GetAll()
	Insert()
	Delete()
}