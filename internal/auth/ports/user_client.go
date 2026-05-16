package ports

import (
	"context"
	"time"
)

type User struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type UserClient interface {
	Ping(ctx context.Context) error

	CreateUser(ctx context.Context, username, password string) (userID string, createdAt time.Time, err error)
	UserByUsername(ctx context.Context, username string) (User, error)
}
