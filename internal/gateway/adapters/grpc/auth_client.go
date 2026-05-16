package grpcstore

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/zennify/backend/gen/go/auth/v1"
	"github.com/zennify/backend/internal/gateway/core/ports"
	"github.com/zennify/backend/internal/shared/grpcaddr"
)

var _ ports.AuthClient = (*AuthClient)(nil)

type AuthClient struct {
	auth authv1.AuthServiceClient
}

func NewAuthClient(conn grpc.ClientConnInterface) *AuthClient {
	return &AuthClient{auth: authv1.NewAuthServiceClient(conn)}
}

func NewAuthConn(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(grpcaddr.DialTarget(addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (c *AuthClient) Register(ctx context.Context, req *ports.RegisterRequest) (*authv1.RegisterResponse, error) {
	return c.auth.Register(ctx, &authv1.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
}

func (c *AuthClient) Login(ctx context.Context, req *ports.LoginRequest) (*authv1.LoginResponse, error) {
	return c.auth.Login(ctx, &authv1.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
}

func (c *AuthClient) RefreshTokens(ctx context.Context, req *ports.RefreshTokensRequest) (*authv1.RefreshTokensResponse, error) {
	return c.auth.RefreshTokens(ctx, &authv1.RefreshTokensRequest{
		RefreshToken: req.RefreshToken,
	})
}
