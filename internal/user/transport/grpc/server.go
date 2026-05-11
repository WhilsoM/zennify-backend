package grpcapi

import (
	"google.golang.org/grpc"

	userv1 "github.com/zennify/backend/gen/go/user/v1"
	"github.com/zennify/backend/internal/user/app"
)

func Register(reg grpc.ServiceRegistrar, svc *app.Service) {
	userv1.RegisterUserServiceServer(reg, newUserServer(svc))
}
