package domain

import "github.com/google/uuid"

type State = int

const (
	EMPTY State = iota
)

type TGUser struct {
	State   State
	UserID  string
	Clients []uuid.UUID
}
