package errors

import "errors"

var (
	ErrorUserAlreadyExists = errors.New("user with this id already exists")
)
