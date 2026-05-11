package grpcapi

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	authv1 "github.com/zennify/backend/gen/go/auth/v1"
	"github.com/zennify/backend/internal/auth/app"
	"github.com/zennify/backend/internal/auth/ports"
)

type authServer struct {
	authv1.UnimplementedAuthServiceServer
	svc *app.Service
	vld *validator.Validate
}

func newAuthServer(svc *app.Service) *authServer {
	return &authServer{
		svc: svc,
		vld: validator.New(),
	}
}

func (a *authServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	in := ports.RegisterRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	if err := a.vld.Struct(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	userID, createdAt, err := a.svc.Register(ctx, in.Username, in.Password)
	if err != nil {
		return nil, mapAuthErr(err)
	}

	return &authv1.RegisterResponse{
		UserId:    userID,
		Username:  in.Username,
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func (a *authServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	in := ports.LoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	if err := a.vld.Struct(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	access, refresh, err := a.svc.Login(ctx, in.Username, in.Password)
	if err != nil {
		return nil, mapAuthErr(err)
	}

	return &authv1.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (a *authServer) RefreshTokens(ctx context.Context, req *authv1.RefreshTokensRequest) (*authv1.RefreshTokensResponse, error) {
	in := ports.RefreshTokensRequest{
		RefreshToken: req.GetRefreshToken(),
	}
	if err := a.vld.Struct(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation: %v", err)
	}

	access, refresh, err := a.svc.RefreshTokens(ctx, in.RefreshToken)
	if err != nil {
		return nil, mapAuthErr(err)
	}

	return &authv1.RefreshTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func mapAuthErr(err error) error {
	switch {
	case errors.Is(err, ports.ErrUsernameTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, ports.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, ports.ErrInvalidRefreshToken):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Errorf(codes.Internal, "auth: %v", err)
	}
}
