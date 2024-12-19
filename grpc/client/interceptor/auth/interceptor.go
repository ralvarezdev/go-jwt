package auth

import (
	"context"
	gojwtgrpc "github.com/ralvarezdev/go-jwt/grpc"
	gojwtgrpcclientmd "github.com/ralvarezdev/go-jwt/grpc/client/metadata"
	gojwtgrpcmd "github.com/ralvarezdev/go-jwt/grpc/server/metadata"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
	goloadergcloud "github.com/ralvarezdev/go-loader/cloud/gcloud"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Interceptor is the interceptor for the authentication
type Interceptor struct {
	accessToken       string
	grpcInterceptions *map[string]gojwtinterception.Interception
}

// NewInterceptor creates a new authentication interceptor
func NewInterceptor(
	tokenSource *oauth.TokenSource,
	grpcInterceptions *map[string]gojwtinterception.Interception,
) (*Interceptor, error) {
	// Check if the token source is nil
	if tokenSource == nil {
		return nil, goloadergcloud.NilTokenSourceError
	}

	// Get the access token from the token source
	token, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	// Check if the gRPC interceptions is nil
	if grpcInterceptions == nil {
		return nil, gojwtgrpc.NilGRPCInterceptionsError
	}

	return &Interceptor{
		accessToken:       token.AccessToken,
		grpcInterceptions: grpcInterceptions,
	}, nil
}

// Authenticate returns a new unary client interceptor that adds authentication metadata to the context
func (i *Interceptor) Authenticate() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) (err error) {
		// Check if the method should be intercepted
		var ctxMetadata *gojwtgrpcclientmd.CtxMetadata
		interception, ok := (*i.grpcInterceptions)[method]
		if !ok || interception == gojwtinterception.None {
			// Create the unauthenticated context metadata
			ctxMetadata, err = gojwtgrpcclientmd.NewUnauthenticatedCtxMetadata(i.accessToken)
		} else {
			// Get metadata from the context
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				return status.Error(codes.Unauthenticated, gojwtgrpc.MissingMetadataError.Error())
			}

			// Get the raw token from the metadata
			rawToken, err := gojwtgrpcmd.GetAuthorizationTokenFromMetadata(md)
			if err != nil {
				return status.Error(codes.Unauthenticated, err.Error())
			}

			// Create the authenticated context metadata
			ctxMetadata, err = gojwtgrpcclientmd.NewAuthenticatedCtxMetadata(i.accessToken, rawToken)
		}

		// Check if there was an error
		if err != nil {
			return status.Error(codes.Aborted, err.Error())
		}

		// Get the gRPC client context with the metadata
		ctx = gojwtgrpcclientmd.GetCtxWithMetadata(ctxMetadata, ctx)

		// Invoke the original invoker
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
