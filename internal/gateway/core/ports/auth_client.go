package ports

import (
	"context"
	"time"

	authv1 "github.com/zennify/backend/gen/go/auth/v1"
)

type AuthClient interface {
	Register(ctx context.Context, req *RegisterRequest) (*authv1.RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*authv1.LoginResponse, error)
	RefreshTokens(ctx context.Context, req *RefreshTokensRequest) (*authv1.RefreshTokensResponse, error)
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type RegisterResponse struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type RefreshTokensRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=1,max=550"`
}
