package grpcapi

import (
	"errors"

	"github.com/zennify/backend/internal/auth/core/ports"
	"github.com/zennify/backend/internal/shared/grpcerr"
)

func toGRPC(err error) error {
	switch {
	case errors.Is(err, ports.ErrUsernameTaken):
		return grpcerr.AlreadyExists(grpcerr.MsgUsernameTaken)
	case errors.Is(err, ports.ErrInvalidCredentials):
		return grpcerr.Unauthenticated(grpcerr.MsgInvalidCredentials)
	case errors.Is(err, ports.ErrInvalidRefreshToken):
		return grpcerr.Unauthenticated(grpcerr.MsgInvalidRefreshToken)
	default:
		return grpcerr.Internal()
	}
}
