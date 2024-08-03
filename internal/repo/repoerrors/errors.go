package repoerrors

import "errors"

var (
	ErrUserExists = errors.New("already exists")
	ErrNotFound   = errors.New("not found")
)
