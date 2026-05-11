package grpcaddr

import (
	"net"
	"strings"
)

// DialTarget returns a valid gRPC dial target.
// @param addr - address to dial
// @return valid gRPC dial target
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
