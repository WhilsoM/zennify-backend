package grpcerr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ClientError(code codes.Code, message string) error {
	return status.Error(code, message)
}

func Convert(err error) *status.Status {
	return status.Convert(err)
}

func ClientMessage(err error) string {
	return status.Convert(err).Message()
}

func Code(err error) codes.Code {
	return status.Convert(err).Code()
}
