package grpcapi

import (
	"google.golang.org/grpc"

	authv1 "github.com/zennify/backend/gen/go/auth/v1"
	"github.com/zennify/backend/internal/auth/core/services"
)

func Register(reg grpc.ServiceRegistrar, svc *services.Service) {
	authv1.RegisterAuthServiceServer(reg, newAuthServer(svc))
}
