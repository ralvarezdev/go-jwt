package context

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtgrpc "github.com/ralvarezdev/go-jwt/grpc"
)

// SetCtxRawToken sets the raw token to the context
func SetCtxRawToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, gojwtgrpc.AuthorizationMetadataKey, token)
}

// SetCtxTokenClaims sets the token claims to the context
func SetCtxTokenClaims(
	ctx context.Context,
	claims *jwt.MapClaims,
) context.Context {
	return context.WithValue(ctx, gojwt.CtxTokenClaimsKey, claims)
}

// GetCtxRawToken gets the raw token from the context
func GetCtxRawToken(ctx context.Context) (string, error) {
	// Get the raw token from the context
	value := ctx.Value(gojwtgrpc.AuthorizationMetadataKey)
	if value == nil {
		return "", MissingTokenError
	}

	// Check the type of the value
	rawToken, ok := value.(string)
	if !ok {
		return "", UnexpectedTokenTypeError
	}

	return rawToken, nil
}

// GetCtxTokenClaims gets the token claims from the context
func GetCtxTokenClaims(ctx context.Context) (*jwt.MapClaims, error) {
	// Get the claims from the context
	value := ctx.Value(gojwt.CtxTokenClaimsKey)
	if value == nil {
		return nil, MissingTokenClaimsError
	}

	// Check the type of the value
	claims, ok := value.(*jwt.MapClaims)
	if !ok {
		return nil, UnexpectedTokenClaimsTypeError
	}
	return claims, nil
}

// GetCtxTokenClaimsSubject gets the token claims subject from the context
func GetCtxTokenClaimsSubject(ctx context.Context) (string, error) {
	// Get the claims from the context
	claims, err := GetCtxTokenClaims(ctx)
	if err != nil {
		return "", err
	}

	// Get the subject from the claims
	subject, ok := (*claims)[gojwt.SubjectClaim].(string)
	if !ok {
		return "", MissingTokenClaimsSubjectError
	}
	return subject, nil
}

// GetCtxTokenClaimsJwtId gets the token claims JWT ID from the context
func GetCtxTokenClaimsJwtId(ctx context.Context) (string, error) {
	// Get the claims from the context
	claims, err := GetCtxTokenClaims(ctx)
	if err != nil {
		return "", err
	}

	// Get the JWT ID from the claims
	jwtId, ok := (*claims)[gojwt.IdClaim].(string)
	if !ok {
		return "", MissingTokenClaimsIdError
	}
	return jwtId, nil
}
