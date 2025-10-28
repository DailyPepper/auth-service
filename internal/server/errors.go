package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) mapErrorToStatus(err error) error {
	errMsg := err.Error()

	switch {
	case contains(errMsg, "email already exists"):
		return status.Error(codes.AlreadyExists, "Email already registered")
	case contains(errMsg, "validation"):
		return status.Error(codes.InvalidArgument, "Validation failed: "+errMsg)
	case contains(errMsg, "password too weak"):
		return status.Error(codes.InvalidArgument, "Password does not meet requirements")
	default:
		return status.Error(codes.Internal, "Registration failed: "+errMsg)
	}
}

func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
