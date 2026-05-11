package grpcaddr

import "testing"

func TestDialTarget(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{":50051", "127.0.0.1:50051"},
		{"127.0.0.1:50051", "127.0.0.1:50051"},
		{"localhost:50051", "localhost:50051"},
		{"auth:50051", "auth:50051"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := DialTarget(tt.in); got != tt.want {
			t.Errorf("DialTarget(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
