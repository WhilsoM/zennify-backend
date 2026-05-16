package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/zennify/backend/internal/user/core/domain"
	"github.com/zennify/backend/internal/user/core/ports"
)

func (s *Service) CreateUser(ctx context.Context, req ports.CreateUserRequest) (ports.CreateUserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ports.CreateUserResponse{}, fmt.Errorf("user services: hash password: %w", err)
	}

	id := uuid.New()
	createdAt := time.Now().UTC()
	username := req.Username

	err = s.users.Create(ctx, domain.User{
		ID:           id,
		Username:     username,
		PasswordHash: string(passwordHash),
		CreatedAt:    createdAt,
	})
	if err != nil {
		return ports.CreateUserResponse{}, err
	}

	return ports.CreateUserResponse{
		UserID:    id.String(),
		Username:  username,
		CreatedAt: createdAt,
	}, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
