package grpcapi

import (
	"errors"

	"github.com/zennify/backend/internal/shared/grpcerr"
	"github.com/zennify/backend/internal/user/core/ports"
)

func toGRPC(err error) error {
	switch {
	case errors.Is(err, ports.ErrUsernameTaken):
		return grpcerr.AlreadyExists(grpcerr.MsgUsernameTaken)
	case errors.Is(err, ports.ErrUserNotFound):
		return grpcerr.NotFound(grpcerr.MsgUserNotFound)
	default:
		return grpcerr.Internal()
	}
}
