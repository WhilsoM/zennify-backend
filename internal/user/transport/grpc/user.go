package grpcapi

import (
	"context"

	"github.com/go-playground/validator/v10"
	userv1 "github.com/zennify/backend/gen/go/user/v1"
	"github.com/zennify/backend/internal/user/app"
	"github.com/zennify/backend/internal/user/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		return nil, status.Errorf(codes.Internal, "user: %v", err)
	}

	return &userv1.CreateUserResponse{
		UserId:    user.UserID,
		Username:  in.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}
