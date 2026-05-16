package ports

import "errors"

var (
	ErrUsernameTaken       = errors.New("auth: username already taken")
	ErrUserNotFound        = errors.New("auth: user not found")
	ErrInvalidCredentials  = errors.New("auth: invalid credentials")
	ErrInvalidRefreshToken = errors.New("auth: invalid refresh token")
)
