package ports

import "errors"

var (
	ErrUsernameTaken = errors.New("user: username already taken")
	ErrUserNotFound  = errors.New("user: user not found")
)
