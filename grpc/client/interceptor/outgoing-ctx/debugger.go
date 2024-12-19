package outgoing_ctx

import (
	"google.golang.org/grpc"
)

// OutgoingCtx interface
type OutgoingCtx interface {
	PrintOutgoingCtx() grpc.UnaryClientInterceptor
}
