package grpcstore

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/zennify/backend/gen/go/user/v1"
	"github.com/zennify/backend/internal/gateway/ports"
	"github.com/zennify/backend/internal/shared/grpcaddr"
)

var _ ports.UserClient = (*UserClient)(nil)

type UserClient struct {
	user userv1.UserServiceClient
}

func NewUserClient(conn grpc.ClientConnInterface) *UserClient {
	return &UserClient{
		user: userv1.NewUserServiceClient(conn),
	}
}

func NewUserConn(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(grpcaddr.DialTarget(addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (c *UserClient) CreateUser(ctx context.Context, req *ports.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	return c.user.CreateUser(ctx, &userv1.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	})
}
