package grpcapi

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	userv1 "github.com/zennify/backend/gen/go/user/v1"
	"github.com/zennify/backend/internal/auth/core/domain"
	"github.com/zennify/backend/internal/auth/core/ports"
	"github.com/zennify/backend/internal/shared/grpcaddr"
)

var _ ports.UserClient = (*UserClient)(nil)

type UserClient struct {
	conn   *grpc.ClientConn
	client userv1.UserServiceClient
}

func NewUserClient(ctx context.Context, addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(
		grpcaddr.DialTarget(addr),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("auth user grpc client: dial: %w", err)
	}

	client := userv1.NewUserServiceClient(conn)
	if _, err := client.Health(ctx, &userv1.HealthRequest{}); err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			return nil, fmt.Errorf("auth user grpc client: health check failed: %w (close: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("auth user grpc client: health check: %w", err)
	}

	return &UserClient{conn: conn, client: client}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) Ping(ctx context.Context) error {
	_, err := c.client.Health(ctx, &userv1.HealthRequest{})
	if err != nil {
		return fmt.Errorf("auth user grpc client: ping: %w", err)
	}
	return nil
}

func (c *UserClient) CreateUser(ctx context.Context, username, password string) (string, time.Time, error) {
	resp, err := c.client.CreateUser(ctx, &userv1.CreateUserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", time.Time{}, mapUserClientErr(err, "create user")
	}
	return resp.GetUserId(), resp.GetCreatedAt().AsTime(), nil
}

func (c *UserClient) UserByUsername(ctx context.Context, username string) (domain.User, error) {
	resp, err := c.client.GetUserByUsername(ctx, &userv1.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return domain.User{}, mapUserClientErr(err, "get user by username")
	}

	return domain.User{
		ID:           resp.GetUserId(),
		Username:     resp.GetUsername(),
		PasswordHash: resp.GetPasswordHash(),
		CreatedAt:    resp.GetCreatedAt().AsTime(),
	}, nil
}

func mapUserClientErr(err error, action string) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("auth user grpc client: %s: %w", action, err)
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return ports.ErrUsernameTaken
	case codes.NotFound:
		return ports.ErrUserNotFound
	default:
		return fmt.Errorf("auth user grpc client: %s: %w", action, err)
	}
}
