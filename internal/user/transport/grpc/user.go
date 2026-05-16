package grpcapi

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/zennify/backend/gen/go/user/v1"
	"github.com/zennify/backend/internal/shared/grpcerr"
	"github.com/zennify/backend/internal/user/app"
	"github.com/zennify/backend/internal/user/ports"
)

type userServer struct {
	userv1.UnimplementedUserServiceServer
	svc *app.Service
	vld *validator.Validate
}

func newUserServer(svc *app.Service) *userServer {
	return &userServer{
		svc: svc,
		vld: validator.New(),
	}
}

func (u *userServer) Health(ctx context.Context, req *userv1.HealthRequest) (*userv1.HealthResponse, error) {
	return &userv1.HealthResponse{
		Message: "OK",
	}, nil
}

func (u *userServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	in := ports.CreateUserRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	if err := u.vld.Struct(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	user, err := u.svc.CreateUser(ctx, in)
	if err != nil {
		if errors.Is(err, ports.ErrUsernameTaken) {
			return nil, grpcerr.ClientError(codes.AlreadyExists, grpcerr.MsgUsernameTaken)
		}
		return nil, grpcerr.ClientError(codes.Internal, grpcerr.MsgInternal)
	}

	return &userv1.CreateUserResponse{
		UserId:    user.UserID,
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (u *userServer) GetUserByUsername(ctx context.Context, req *userv1.GetUserByUsernameRequest) (*userv1.GetUserByUsernameResponse, error) {
	username := req.GetUsername()
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	user, err := u.svc.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return nil, grpcerr.ClientError(codes.NotFound, grpcerr.MsgUserNotFound)
		}
		return nil, grpcerr.ClientError(codes.Internal, grpcerr.MsgInternal)
	}

	return &userv1.GetUserByUsernameResponse{
		UserId:       user.ID.String(),
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		CreatedAt:    timestamppb.New(user.CreatedAt),
	}, nil
}
