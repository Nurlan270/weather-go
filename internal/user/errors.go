package user

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("password is invalid")
	ErrSessionExpired    = errors.New("session expired")
)
