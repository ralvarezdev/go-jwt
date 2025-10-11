package context

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
//   - key: The key to set the token under
//   - token: The raw token to set
//
// Returns:
//
//   - context.Context: The context with the raw token set
func SetCtxToken(ctx context.Context, key, token string) context.Context {
	return context.WithValue(ctx, key, token)
}

// SetCtxRefreshToken sets the raw refresh token to the context
//
// Parameters:
//
//   - ctx: The context to set the raw refresh token to
//   - token: The raw refresh token to set
//
// Returns:
//
//   - context.Context: The context with the raw refresh token set
func SetCtxRefreshToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, gojwt.CtxRefreshTokenKey, token)
}

// SetCtxAccessToken sets the raw access token to the context
//
// Parameters:
//
//   - ctx: The context to set the raw access token to
//   - token: The raw access token to set
//
// Returns:
//
//   - context.Context: The context with the raw access token set
func SetCtxAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, gojwt.CtxAccessTokenKey, token)
}

// SetCtxTokenClaims sets the token claims to the context
//
// Parameters:
//
//   - ctx: The context to set the token claims to
//   - key: The key to set the claims under
//   - claims: The token claims to set
//
// Returns:
//
//   - context.Context: The context with the token claims set
func SetCtxTokenClaims(
	ctx context.Context,
	key string,
	claims jwt.MapClaims,
) context.Context {
	return context.WithValue(ctx, key, claims)
}

// SetCtxRefreshTokenClaims sets the refresh token claims to the context
//
// Parameters:
//
//   - ctx: The context to set the refresh token claims to
//   - claims: The refresh token claims to set
//
// Returns:
//
//   - context.Context: The context with the refresh token claims set
func SetCtxRefreshTokenClaims(
	ctx context.Context,
	claims jwt.MapClaims,
) context.Context {
	return context.WithValue(ctx, gojwt.CtxRefreshTokenClaimsKey, claims)
}

// SetCtxAccessTokenClaims sets the access token claims to the context
//
// Parameters:
//
//   - ctx: The context to set the access token claims to
//   - claims: The access token claims to set
//
// Returns:
//
//   - context.Context: The context with the access token claims set
func SetCtxAccessTokenClaims(
	ctx context.Context,
	claims jwt.MapClaims,
) context.Context {
	return context.WithValue(ctx, gojwt.CtxAccessTokenClaimsKey, claims)
}

// GetCtxToken gets the raw token from the context
//
// Parameters:
//
//   - ctx: The context to get the raw token from
//   - key: The key to get the token from
//
// Returns:
//
//   - string: The raw token
//   - error: An error if the raw token is not found or is of an unexpected type
func GetCtxToken(ctx context.Context, key string) (string, error) {
	// Get the raw token from the context
	value := ctx.Value(key)
	if value == nil {
		return "", ErrMissingToken
	}

	// Check the type of the value
	rawToken, ok := value.(string)
	if !ok {
		return "", ErrUnexpectedTokenType
	}

	return rawToken, nil
}

// GetCtxRefreshToken gets the raw refresh token from the context
//
// Parameters:
//
//   - ctx: The context to get the raw refresh token from
//
// Returns:
//
//   - string: The raw refresh token
//   - error: An error if the raw refresh token is not found or is of an unexpected type
func GetCtxRefreshToken(ctx context.Context) (string, error) {
	return GetCtxToken(ctx, gojwt.CtxRefreshTokenKey)
}

// GetCtxAccessToken gets the raw access token from the context
//
// Parameters:
//
//   - ctx: The context to get the raw access token from
//
// Returns:
//
//   - string: The raw access token
//   - error: An error if the raw access token is not found or is of an unexpected type
func GetCtxAccessToken(ctx context.Context) (string, error) {
	return GetCtxToken(ctx, gojwt.CtxAccessTokenKey)
}

// GetCtxTokenClaims gets the token claims from the context
//
// Parameters:
//
//   - ctx: The context to get the token claims from
//   - key: The key to get the claims from
//
// Returns:
//
//   - jwt.MapClaims: The token claims
//   - error: An error if the token claims are not found or are of an unexpected type
func GetCtxTokenClaims(ctx context.Context, key string) (jwt.MapClaims, error) {
	// Get the claims from the context
	value := ctx.Value(key)
	if value == nil {
		return nil, ErrMissingTokenClaims
	}

	// Check the type of the value
	claims, ok := value.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnexpectedTokenClaimsType
	}
	return claims, nil
}

// GetCtxRefreshTokenClaims gets the refresh token claims from the context
//
// Parameters:
//
//   - ctx: The context to get the refresh token claims from
//
// Returns:
//
//   - jwt.MapClaims: The refresh token claims
//   - error: An error if the refresh token claims are not found or are of an unexpected type
func GetCtxRefreshTokenClaims(ctx context.Context) (jwt.MapClaims, error) {
	return GetCtxTokenClaims(ctx, gojwt.CtxRefreshTokenClaimsKey)
}

// GetCtxAccessTokenClaims gets the access token claims from the context
//
// Parameters:
//
//   - ctx: The context to get the access token claims from
//
// Returns:
//
//   - jwt.MapClaims: The access token claims
//   - error: An error if the access token claims are not found or are of an unexpected type
func GetCtxAccessTokenClaims(ctx context.Context) (jwt.MapClaims, error) {
	return GetCtxTokenClaims(ctx, gojwt.CtxAccessTokenClaimsKey)
}

// GetCtxTokenClaimsSubject gets the token claims subject from the context
//
// Parameters:
//
//   - ctx: The context to get the token claims subject from
//   - key: The key to get the claims from
//
// Returns:
//
//   - string: The token claims subject
//   - error: An error if the token claims subject is not found or is of an unexpected type
func GetCtxTokenClaimsSubject(ctx context.Context, key string) (string, error) {
	// Get the claims from the context
	claims, err := GetCtxTokenClaims(ctx, key)
	if err != nil {
		return "", err
	}

	// Get the subject from the claims
	subject, ok := claims[gojwt.SubjectClaim].(string)
	if !ok {
		return "", ErrMissingTokenClaimsSubject
	}
	return subject, nil
}

// GetCtxRefreshTokenClaimsSubject gets the refresh token claims subject from the context
//
// Parameters:
//
//   - ctx: The context to get the refresh token claims subject from
//
// Returns:
//
//   - string: The refresh token claims subject
//   - error: An error if the refresh token claims subject is not found or is of an unexpected type
func GetCtxRefreshTokenClaimsSubject(ctx context.Context) (string, error) {
	return GetCtxTokenClaimsSubject(ctx, gojwt.CtxRefreshTokenClaimsKey)
}

// GetCtxAccessTokenClaimsSubject gets the access token claims subject from the context
//
// Parameters:
//
//   - ctx: The context to get the access token claims subject from
//
// Returns:
//
//   - string: The access token claims subject
//   - error: An error if the access token claims subject is not found or is of an unexpected type
func GetCtxAccessTokenClaimsSubject(ctx context.Context) (string, error) {
	return GetCtxTokenClaimsSubject(ctx, gojwt.CtxAccessTokenClaimsKey)
}

// GetCtxTokenClaimsJwtID gets the token claims JWT ID from the context
//
// Parameters:
//
//   - ctx: The context to get the token claims JWT ID from
//   - key: The key to get the claims from
//
// Returns:
//
//   - string: The token claims JWT ID
//   - error: An error if the token claims JWT ID is not found or is of an unexpected type
func GetCtxTokenClaimsJwtID(ctx context.Context, key string) (string, error) {
	// Get the claims from the context
	claims, err := GetCtxTokenClaims(ctx, key)
	if err != nil {
		return "", err
	}

	// Get the JWT ID from the claims
	jwtID, ok := claims[gojwt.IDClaim].(string)
	if !ok {
		return "", ErrMissingTokenClaimsID
	}
	return jwtID, nil
}

// GetCtxRefreshTokenClaimsJwtID gets the refresh token claims JWT ID from the context
//
// Parameters:
//
//   - ctx: The context to get the refresh token claims JWT ID from
//
// Returns:
//
//   - string: The refresh token claims JWT ID
//   - error: An error if the refresh token claims JWT ID is not found or is of an unexpected type
func GetCtxRefreshTokenClaimsJwtID(ctx context.Context) (string, error) {
	return GetCtxTokenClaimsJwtID(ctx, gojwt.CtxRefreshTokenClaimsKey)
}

// GetCtxAccessTokenClaimsJwtID gets the access token claims JWT ID from the context
//
// Parameters:
//
//   - ctx: The context to get the access token claims JWT ID from
//
// Returns:
//
//   - string: The access token claims JWT ID
//   - error: An error if the access token claims JWT ID is not found or is of an unexpected type
func GetCtxAccessTokenClaimsJwtID(ctx context.Context) (string, error) {
	return GetCtxTokenClaimsJwtID(ctx, gojwt.CtxAccessTokenClaimsKey)
}

// ClearCtxTokenClaims clears the token claims from the context
//
// Parameters:
//
//   - ctx: The context to clear the token claims from
//   - key: The key to clear the claims from
//
// Returns:
//
//   - context.Context: The context with the token claims cleared
//   - error: An error if the context is nil
func ClearCtxTokenClaims(ctx context.Context, key string) (
	context.Context,
	error,
) {
	return context.WithValue(ctx, key, nil), nil
}

// ClearCtxRefreshTokenClaims clears the refresh token claims from the context
//
// Parameters:
//
//   - ctx: The context to clear the refresh token claims from
//
// Returns:
//
//   - context.Context: The context with the refresh token claims cleared
//   - error: An error if the context is nil
func ClearCtxRefreshTokenClaims(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, gojwt.CtxRefreshTokenClaimsKey, nil), nil
}

// ClearCtxAccessTokenClaims clears the access token claims from the context
//
// Parameters:
//
//   - ctx: The context to clear the access token claims from
//
// Returns:
//
//   - context.Context: The context with the access token claims cleared
//   - error: An error if the context is nil
func ClearCtxAccessTokenClaims(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, gojwt.CtxAccessTokenClaimsKey, nil), nil
}

// ClearCtxTokensClaims clears all tokens claims from the context
//
// Parameters:
//
//   - ctx: The context to clear the tokens claims from
//
// Returns:
//
//   - context.Context: The context with the tokens claims cleared
//   - error: An error if the context is nil
func ClearCtxTokensClaims(ctx context.Context) (context.Context, error) {
	var err error
	ctx, err = ClearCtxRefreshTokenClaims(ctx)
	if err != nil {
		return ctx, err
	}
	return ClearCtxAccessTokenClaims(ctx)
}

// ClearCtxToken clears the raw token from the context
//
// Parameters:
//
//   - ctx: The context to clear the raw token from
//   - key: The key to clear the token from
//
// Returns:
//
//   - context.Context: The context with the raw token cleared
//   - error: An error if the context is nil
func ClearCtxToken(ctx context.Context, key string) (context.Context, error) {
	return context.WithValue(ctx, key, nil), nil
}

// ClearCtxRefreshToken clears the raw refresh token from the context
//
// Parameters:
//
//   - ctx: The context to clear the raw refresh token from
//
// Returns:
//
//   - context.Context: The context with the raw refresh token cleared
//   - error: An error if the context is nil
func ClearCtxRefreshToken(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, gojwt.CtxRefreshTokenKey, nil), nil
}

// ClearCtxAccessToken clears the raw access token from the context
//
// Parameters:
//
//   - ctx: The context to clear the raw access token from
//
// Returns:
//
//   - context.Context: The context with the raw access token cleared
//   - error: An error if the context is nil
func ClearCtxAccessToken(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, gojwt.CtxAccessTokenKey, nil), nil
}

// ClearCtxTokens clears all tokens from the context
//
// Parameters:
//
//   - ctx: The context to clear the tokens from
//
// Returns:
//
//   - context.Context: The context with the tokens cleared
//   - error: An error if the context is nil
func ClearCtxTokens(ctx context.Context) (context.Context, error) {
	var err error
	ctx, err = ClearCtxRefreshToken(ctx)
	if err != nil {
		return ctx, err
	}
	return ClearCtxAccessToken(ctx)
}

// ClearCtxAll clears all tokens and claims from the context
//
// Parameters:
//
//   - ctx: The context to clear the tokens and claims from
//
// Returns:
//
//   - context.Context: The context with the tokens and claims cleared
//   - error: An error if the context is nil
func ClearCtxAll(ctx context.Context) (context.Context, error) {
	var err error
	ctx, err = ClearCtxTokens(ctx)
	if err != nil {
		return ctx, err
	}
	return ClearCtxTokensClaims(ctx)
}
