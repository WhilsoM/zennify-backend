// Package grpcerr builds gRPC status errors with stable client messages (Msg*).
//
// Use only in gRPC adapters (internal/*/adapters/grpc).
// HTTP gateway maps upstream gRPC errors in adapters/http/errors.go.
package grpcerr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InvalidRequest() error {
	return status.Error(codes.InvalidArgument, MsgInvalidRequest)
}

func Internal() error {
	return status.Error(codes.Internal, MsgInternal)
}

func NotFound(msg string) error {
	return status.Error(codes.NotFound, msg)
}

func AlreadyExists(msg string) error {
	return status.Error(codes.AlreadyExists, msg)
}

func Unauthenticated(msg string) error {
	return status.Error(codes.Unauthenticated, msg)
}
