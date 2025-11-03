package grpc

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
)

// SetCtxToken sets the raw token to the context
//
// Parameters:
//
//   - ctx: The context to set the raw token to
//   - token: The raw token to set
//
// Returns:
//
//   - context.Context: The context with the raw token set
func SetCtxToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, gojwt.CtxTokenKey, token)
}

// SetCtxTokenClaims sets the token claims to the context
//
// Parameters:
//
//   - ctx: The context to set the token claims to
//   - claims: The token claims to set
//
// Returns:
//
//   - context.Context: The context with the token claims set
func SetCtxTokenClaims(
	ctx context.Context,
	claims jwt.MapClaims,
) context.Context {
	return context.WithValue(ctx, gojwt.CtxTokenClaimsKey, claims)
}

// GetCtxToken gets the raw token from the context
//
// Parameters:
//
//   - ctx: The context to get the raw token from
//
// Returns:
//
//   - string: The raw token
//   - error: An error if the raw token is not found or is of an unexpected type
func GetCtxToken(ctx context.Context) (string, error) {
	// Get the raw token from the context
	value := ctx.Value(gojwt.CtxTokenKey)
	if value == nil {
		return "", gojwt.ErrMissingTokenInContext
	}

	// Check the type of the value
	rawToken, ok := value.(string)
	if !ok {
		return "", gojwt.ErrUnexpectedTokenTypeInContext
	}

	return rawToken, nil
}

// GetCtxTokenClaims gets the token claims from the context
//
// Parameters:
//
//   - ctx: The context to get the token claims from
//
// Returns:
//
//   - jwt.MapClaims: The token claims
//   - error: An error if the token claims are not found or are of an unexpected type
func GetCtxTokenClaims(ctx context.Context) (jwt.MapClaims, error) {
	// Get the claims from the context
	value := ctx.Value(gojwt.CtxTokenClaimsKey)
	if value == nil {
		return nil, gojwt.ErrMissingTokenClaimsInContext
	}

	// Check the type of the value
	claims, ok := value.(jwt.MapClaims)
	if !ok {
		return nil, gojwt.ErrUnexpectedTokenClaimsTypeInContext
	}
	return claims, nil
}

// GetCtxTokenClaimsSubject gets the token claims subject from the context
//
// Parameters:
//
//   - ctx: The context to get the token claims subject from
//
// Returns:
//
//   - string: The token claims subject
//   - error: An error if the token claims subject is not found or is of an unexpected type
func GetCtxTokenClaimsSubject(ctx context.Context) (string, error) {
	// Get the claims from the context
	claims, err := GetCtxTokenClaims(ctx)
	if err != nil {
		return "", err
	}

	// Get the subject from the claims
	subject, ok := claims[gojwt.SubjectClaim].(string)
	if !ok {
		return "", gojwt.ErrMissingTokenClaimsSubject
	}
	return subject, nil
}

// GetCtxTokenClaimsJwtID gets the token claims JWT ID from the context
//
// Parameters:
//
//   - ctx: The context to get the token claims JWT ID from
//
// Returns:
//
//   - string: The token claims JWT ID
//   - error: An error if the token claims JWT ID is not found or is of an unexpected type
func GetCtxTokenClaimsJwtID(ctx context.Context) (string, error) {
	// Get the claims from the context
	claims, err := GetCtxTokenClaims(ctx)
	if err != nil {
		return "", err
	}

	// Get the JWT ID from the claims
	jwtID, ok := claims[gojwt.IDClaim].(string)
	if !ok {
		return "", gojwt.ErrMissingTokenClaimsID
	}
	return jwtID, nil
}
