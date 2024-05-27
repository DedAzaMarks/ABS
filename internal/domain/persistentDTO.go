package domain

import "github.com/google/uuid"

type DeviceDTO struct {
	ID   uuid.UUID
	Name string
}

type UserDTO struct {
	ID         int64
	SessionKey string
	State      string
	Devices    []DeviceDTO
}
