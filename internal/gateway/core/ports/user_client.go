package ports

import (
	"context"
	"time"
)

type UserClient interface {
	GetUserByID(ctx context.Context, userID string) (*UserProfile, error)
}

type UserProfile struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
