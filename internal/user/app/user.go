package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zennify/backend/internal/user/ports"
)

func (s *Service) CreateUser(ctx context.Context, req ports.CreateUserRequest) (ports.CreateUserResponse, error) {
	return ports.CreateUserResponse{
		UserID:    uuid.New().String(),
		Username:  req.Username,
		CreatedAt: time.Now(),
	}, nil
}
