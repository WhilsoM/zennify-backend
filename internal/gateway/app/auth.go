package app

import (
	"context"

	authv1 "github.com/zennify/backend/gen/go/auth/v1"
	"github.com/zennify/backend/internal/gateway/ports"
)

func (s *Service) Register(ctx context.Context, req *ports.RegisterRequest) (*ports.RegisterResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
	defer cancel()
	out, err := s.breaker.Execute(func() (any, error) {
		return s.auth.Register(callCtx, req)
	})
	if err != nil {
		return nil, err
	}
	pb := out.(*authv1.RegisterResponse)
	resp := &ports.RegisterResponse{
		UserID:   pb.GetUserId(),
		Username: pb.GetUsername(),
	}
	if ts := pb.GetCreatedAt(); ts != nil {
		resp.CreatedAt = ts.AsTime()
	}
	return resp, nil
}

func (s *Service) Login(ctx context.Context, req *ports.LoginRequest) (*authv1.LoginResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
	defer cancel()
	out, err := s.breaker.Execute(func() (any, error) {
		return s.auth.Login(callCtx, req)
	})
	if err != nil {
		return nil, err
	}
	return out.(*authv1.LoginResponse), nil
}

func (s *Service) RefreshTokens(ctx context.Context, req *ports.RefreshTokensRequest) (*authv1.RefreshTokensResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, s.requestTimeout)
	defer cancel()
	out, err := s.breaker.Execute(func() (any, error) {
		return s.auth.RefreshTokens(callCtx, req)
	})
	if err != nil {
		return nil, err
	}
	return out.(*authv1.RefreshTokensResponse), nil
}
