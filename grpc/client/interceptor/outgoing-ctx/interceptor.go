package outgoing_ctx

import (
	"context"
	gologger "github.com/ralvarezdev/go-logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Interceptor is the interceptor for the outgoing context
type Interceptor struct {
	logger *Logger
}

// NewInterceptor creates a new interceptor for the outgoing context
func NewInterceptor(logger *Logger) (*Interceptor, error) {
	// Check if the logger is nil
	if logger == nil {
		return nil, gologger.NilLoggerError
	}

	return &Interceptor{
		logger: logger,
	}, nil
}

// PrintOutgoingCtx prints the outgoing context
func (i *Interceptor) PrintOutgoingCtx() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Get the outgoing context
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			return status.Error(codes.Internal, FailedToGetOutgoingContextError.Error())
		}

		// Print the metadata
		for key, values := range md {
			for _, value := range values {
				i.logger.LogKeyValue(key, value)
			}
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
