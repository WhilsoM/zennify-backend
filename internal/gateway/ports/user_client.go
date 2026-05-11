package ports

import (
	"context"
	"time"

	userv1 "github.com/zennify/backend/gen/go/user/v1"
)

type UserClient interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*userv1.CreateUserResponse, error)
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type CreateUserResponse struct {
	UserID    string
	Username  string
	CreatedAt time.Time
}
