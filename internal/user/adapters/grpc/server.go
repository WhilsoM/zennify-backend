package grpcapi

import (
	"google.golang.org/grpc"

	userv1 "github.com/zennify/backend/gen/go/user/v1"
	"github.com/zennify/backend/internal/user/core/services"
)

func Register(reg grpc.ServiceRegistrar, svc *services.Service) {
	userv1.RegisterUserServiceServer(reg, newUserServer(svc))
}
