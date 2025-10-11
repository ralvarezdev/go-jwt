package context

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
)

// SetCtxTokenClaims sets the token claims in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - key: The key to set the claims under
//   - claims: The token claims to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the token claims set in the context
func SetCtxTokenClaims(
	r *http.Request,
	key string,
	claims jwt.MapClaims,
) *http.Request {
	ctx := context.WithValue(r.Context(), key, claims)
	return r.WithContext(ctx)
}

// SetCtxRefreshTokenClaims sets the refresh token claims in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - claims: The refresh token claims to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the refresh token claims set in the context
//   - error: An error if the request is nil
func SetCtxRefreshTokenClaims(
	r *http.Request,
	claims jwt.MapClaims,
) (*http.Request, error) {
	return SetCtxTokenClaims(r, gojwt.CtxRefreshTokenClaimsKey, claims), nil
}

// SetCtxAccessTokenClaims sets the access token claims in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - claims: The access token claims to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the access token claims set in the context
//   - error: An error if the request is nil
func SetCtxAccessTokenClaims(
	r *http.Request,
	claims jwt.MapClaims,
) (*http.Request, error) {
	return SetCtxTokenClaims(r, gojwt.CtxAccessTokenClaimsKey, claims), nil
}

// GetCtxTokenClaims tries to get the token claims from the context
//
// Parameters:
//
//   - r: The HTTP request
//   - key: The key to get the claims from
//
// Returns:
//
//   - jwt.MapClaims: The token claims from the context
//   - error: An error if the token claims are not found or of an unexpected type
func GetCtxTokenClaims(r *http.Request, key string) (jwt.MapClaims, error) {
	// Get the token claims from the context
	value := r.Context().Value(key)
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

// GetCtxRefreshTokenClaims tries to get the refresh token claims from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - jwt.MapClaims: The refresh token claims from the context
//   - error: An error if the refresh token claims are not found or of an unexpected type
func GetCtxRefreshTokenClaims(r *http.Request) (jwt.MapClaims, error) {
	return GetCtxTokenClaims(r, gojwt.CtxRefreshTokenClaimsKey)
}

// GetCtxAccessTokenClaims tries to get the access token claims from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - jwt.MapClaims: The access token claims from the context
//   - error: An error if the access token claims are not found or of an unexpected type
func GetCtxAccessTokenClaims(r *http.Request) (jwt.MapClaims, error) {
	return GetCtxTokenClaims(r, gojwt.CtxAccessTokenClaimsKey)
}

// SetCtxToken sets the raw token in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - key: The key to set the token under
//   - token: The raw token to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the raw token set in the context
//   - error: An error if the token is empty
func SetCtxToken(r *http.Request, key, token string) (*http.Request, error) {
	// Check if the token is empty
	if token == "" {
		return nil, gojwt.ErrEmptyToken
	}

	ctx := context.WithValue(r.Context(), key, token)
	return r.WithContext(ctx), nil
}

// SetCtxRefreshToken sets the raw refresh token in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - token: The raw refresh token to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the raw refresh token set in the context
//   - error: An error if the token is empty
func SetCtxRefreshToken(r *http.Request, token string) (*http.Request, error) {
	return SetCtxToken(r, gojwt.CtxRefreshTokenKey, token)
}

// SetCtxAccessToken sets the raw access token in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - token: The raw access token to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the raw access token set in the context
//   - error: An error if the token is empty
func SetCtxAccessToken(r *http.Request, token string) (*http.Request, error) {
	return SetCtxToken(r, gojwt.CtxAccessTokenKey, token)
}

// GetCtxToken tries to get the raw token from the context
//
// Parameters:
//
//   - r: The HTTP request
//   - key: The key to get the token from
//
// Returns:
//
//   - string: The raw token from the context
//   - error: An error if the token is not found or of an unexpected type
func GetCtxToken(r *http.Request, key string) (string, error) {
	// Get the token from the context
	value := r.Context().Value(key)
	if value == nil {
		return "", gojwt.ErrMissingTokenInContext
	}

	// Check the type of the value
	token, ok := value.(string)
	if !ok {
		return "", gojwt.ErrUnexpectedTokenTypeInContext
	}

	return token, nil
}

// GetCtxRefreshToken tries to get the raw refresh token from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - string: The raw refresh token from the context
//   - error: An error if the refresh token is not found or of an unexpected type
func GetCtxRefreshToken(r *http.Request) (string, error) {
	return GetCtxToken(r, gojwt.CtxRefreshTokenKey)
}

// GetCtxAccessToken tries to get the raw access token from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - string: The raw access token from the context
//   - error: An error if the access token is not found or of an unexpected type
func GetCtxAccessToken(r *http.Request) (string, error) {
	return GetCtxToken(r, gojwt.CtxAccessTokenKey)
}
