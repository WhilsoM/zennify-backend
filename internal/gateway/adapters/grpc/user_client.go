package grpcstore

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/zennify/backend/gen/go/user/v1"
	"github.com/zennify/backend/internal/gateway/core/ports"
	"github.com/zennify/backend/internal/shared/grpcaddr"
)

var _ ports.UserClient = (*UserClient)(nil)

type UserClient struct {
	user userv1.UserServiceClient
}

func NewUserClient(conn grpc.ClientConnInterface) *UserClient {
	return &UserClient{user: userv1.NewUserServiceClient(conn)}
}

func NewUserConn(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(grpcaddr.DialTarget(addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (c *UserClient) GetUserByID(ctx context.Context, userID string) (*ports.UserProfile, error) {
	req := &userv1.GetUserByIDRequest{UserId: userID}

	resp, err := c.user.GetUserByID(ctx, req)
	if err != nil {
		return nil, err
	}

	profile := &ports.UserProfile{
		UserID:   resp.GetUserId(),
		Username: resp.GetUsername(),
	}
	if ts := resp.GetCreatedAt(); ts != nil {
		profile.CreatedAt = ts.AsTime()
	}
	return profile, nil
}
