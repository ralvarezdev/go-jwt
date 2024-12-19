package auth

import (
	"context"
	"errors"
	gojwtgrpc "github.com/ralvarezdev/go-jwt/grpc"
	gojwtserver "github.com/ralvarezdev/go-jwt/grpc/server"
	gojwtserverctx "github.com/ralvarezdev/go-jwt/grpc/server/context"
	gojwtservermd "github.com/ralvarezdev/go-jwt/grpc/server/metadata"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
	gojwtvalidator "github.com/ralvarezdev/go-jwt/token/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Interceptor is the interceptor for the authentication
type Interceptor struct {
	validator         gojwtvalidator.Validator
	grpcInterceptions *map[string]gojwtinterception.Interception
}

// NewInterceptor creates a new authentication interceptor
func NewInterceptor(
	validator gojwtvalidator.Validator,
	grpcInterceptions *map[string]gojwtinterception.Interception,
) (*Interceptor, error) {
	// Check if either the validator or the gRPC interceptions is nil
	if validator == nil {
		return nil, gojwtvalidator.NilValidatorError
	}
	if grpcInterceptions == nil {
		return nil, gojwtgrpc.NilGRPCInterceptionsError
	}

	return &Interceptor{
		validator:         validator,
		grpcInterceptions: grpcInterceptions,
	}, nil
}

// Authenticate returns the authentication interceptor
func (i *Interceptor) Authenticate() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if the method should be intercepted
		interception, ok := (*i.grpcInterceptions)[info.FullMethod]
		if !ok || interception == gojwtinterception.None {
			return handler(ctx, req)
		}

		// Get metadata from the context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, gojwtgrpc.MissingMetadataError.Error())
		}

		// Get the raw token from the metadata
		rawToken, err := gojwtservermd.GetAuthorizationTokenFromMetadata(md)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		// Validate the token and get the validated claims
		claims, err := i.validator.GetValidatedClaims(rawToken, interception)
		if err != nil {
			if errors.Is(err, gojwtvalidator.NilClaimsError) {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}

			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, status.Error(codes.Unauthenticated, gojwtserver.TokenHasExpiredError.Error())
			}

			return nil, status.Error(codes.Internal, gogrpc.InternalServerError.Error())
		}

		// Set the raw token and token claims to the context
		ctx = gojwtserverctx.SetCtxRawToken(ctx, rawToken)
		ctx = gojwtserverctx.SetCtxTokenClaims(ctx, claims)

		return handler(ctx, req)
	}
}
