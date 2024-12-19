package auth

import (
	"google.golang.org/grpc"
)

// Authentication interface
type Authentication interface {
	Authenticate() grpc.UnaryServerInterceptor
}
