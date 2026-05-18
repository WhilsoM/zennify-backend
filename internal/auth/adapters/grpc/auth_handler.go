package grpcapi

import (
	"context"

	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/types/known/timestamppb"

	authv1 "github.com/zennify/backend/gen/go/auth/v1"
	"github.com/zennify/backend/internal/auth/core/ports"
	"github.com/zennify/backend/internal/auth/core/services"
	"github.com/zennify/backend/internal/shared/grpcerr"
)

type authServer struct {
	authv1.UnimplementedAuthServiceServer
	svc *services.Service
	vld *validator.Validate
}

func newAuthServer(svc *services.Service) *authServer {
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
		return nil, grpcerr.InvalidRequest()
	}

	userID, createdAt, err := a.svc.Register(ctx, in)
	if err != nil {
		return nil, toGRPC(err)
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
		return nil, grpcerr.InvalidRequest()
	}

	access, refresh, err := a.svc.Login(ctx, in)
	if err != nil {
		return nil, toGRPC(err)
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
		return nil, grpcerr.InvalidRequest()
	}

	access, refresh, err := a.svc.RefreshTokens(ctx, in)
	if err != nil {
		return nil, toGRPC(err)
	}

	return &authv1.RefreshTokensResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
