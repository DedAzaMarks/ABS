package errors

import "errors"

var (
	ErrorUserNotFound      = errors.New("user not found")
	ErrorDeviceNotFound    = errors.New("device not found")
	ErrorUserAlreadyExists = errors.New("user with this id already exists")
)

var (
	NotAFilm = errors.New("not a film")
)

type ErrFromUser struct {
	error
	UserID int64
}
