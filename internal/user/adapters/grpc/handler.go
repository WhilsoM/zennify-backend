package grpcapi

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/zennify/backend/gen/go/user/v1"
	"github.com/zennify/backend/internal/shared/grpcerr"
	"github.com/zennify/backend/internal/user/core/ports"
	"github.com/zennify/backend/internal/user/core/services"
)

type userServer struct {
	userv1.UnimplementedUserServiceServer
	svc *services.Service
	vld *validator.Validate
}

func newUserServer(svc *services.Service) *userServer {
	return &userServer{
		svc: svc,
		vld: validator.New(),
	}
}

func (u *userServer) Health(ctx context.Context, req *userv1.HealthRequest) (*userv1.HealthResponse, error) {
	_ = ctx
	_ = req
	return &userv1.HealthResponse{Message: "OK"}, nil
}

func (u *userServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	in := ports.CreateUserRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	if err := u.vld.Struct(in); err != nil {
		return nil, grpcerr.ClientError(codes.InvalidArgument, grpcerr.MsgInvalidRequest)
	}

	user, err := u.svc.CreateUser(ctx, in)
	if err != nil {
		return nil, mapUserErr(err)
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
		return nil, grpcerr.ClientError(codes.InvalidArgument, grpcerr.MsgInvalidRequest)
	}

	user, err := u.svc.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, mapUserErr(err)
	}

	return &userv1.GetUserByUsernameResponse{
		UserId:       user.ID.String(),
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		CreatedAt:    timestamppb.New(user.CreatedAt),
	}, nil
}

func mapUserErr(err error) error {
	switch {
	case errors.Is(err, ports.ErrUsernameTaken):
		return grpcerr.ClientError(codes.AlreadyExists, grpcerr.MsgUsernameTaken)
	case errors.Is(err, ports.ErrUserNotFound):
		return grpcerr.ClientError(codes.NotFound, grpcerr.MsgUserNotFound)
	default:
		return grpcerr.ClientError(codes.Internal, grpcerr.MsgInternal)
	}
}
