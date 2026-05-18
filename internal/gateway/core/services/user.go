package services

import (
	"context"

	"github.com/zennify/backend/internal/gateway/core/ports"
)

func (s *Service) GetUserProfile(ctx context.Context, userID string) (*ports.UserProfile, error) {
	callCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
	defer cancel()

	out, err := s.breaker.Execute(func() (any, error) {
		return s.users.GetUserByID(callCtx, userID)
	})
	if err != nil {
		return nil, err
	}
	return out.(*ports.UserProfile), nil
}
