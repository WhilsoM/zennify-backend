package services

import "github.com/zennify/backend/internal/user/core/ports"

type Service struct {
	users ports.UserRepository
}

func NewService(users ports.UserRepository) *Service {
	return &Service{users: users}
}
