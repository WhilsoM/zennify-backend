// Package grpcaddr normalizes host:port strings for gRPC client dial targets.
//
// Service config uses ":50052" (all interfaces on the server). grpc.Dial
// requires an explicit host, so callers use DialTarget before opening a connection.
package grpcaddr

import (
	"net"
	"strings"
)

// DialTarget returns addr in a form valid for grpc.NewClient and grpc.Dial.
//
// Leading and trailing whitespace is trimmed. An empty string is returned unchanged.
//
// Addresses with a missing host are normalized to loopback:
//   - ":50052" -> "127.0.0.1:50052"
//   - "localhost:8080" -> unchanged
//
// Values that are not host:port (SplitHostPort fails and addr does not start with ":")
// are returned as-is so callers can pass DNS names or resolver schemes unchanged.
func DialTarget(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return addr
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		if strings.HasPrefix(addr, ":") {
			return "127.0.0.1" + addr
		}
		return addr
	}
	if host == "" {
		return net.JoinHostPort("127.0.0.1", port)
	}
	return addr
}
