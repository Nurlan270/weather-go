package location

import "errors"

var (
	ErrNoResults             = errors.New("no results found")
	ErrLocationAlreadyExists = errors.New("location already exists")
)
