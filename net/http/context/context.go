package context

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtnethttp "github.com/ralvarezdev/go-jwt/net/http"
)

// SetCtxTokenClaims sets the token claims in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - claims: The token claims to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the token claims set in the context
func SetCtxTokenClaims(
	r *http.Request,
	claims jwt.MapClaims,
) *http.Request {
	ctx := context.WithValue(r.Context(), gojwtnethttp.AuthorizationKey, claims)
	return r.WithContext(ctx)
}

// GetCtxTokenClaims tries to get the token claims from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - jwt.MapClaims: The token claims from the context
//   - error: An error if the token claims are not found or of an unexpected type
func GetCtxTokenClaims(r *http.Request) (jwt.MapClaims, error) {
	// Get the token claims from the context
	value := r.Context().Value(gojwtnethttp.AuthorizationKey)
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

// SetCtxToken sets the raw token in the context
//
// Parameters:
//
//   - r: The HTTP request
//   - token: The raw token to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the raw token set in the context
//   - error: An error if the token is empty
func SetCtxToken(r *http.Request, token string) (*http.Request, error) {
	// Check if the token is empty
	if token == "" {
		return nil, gojwt.ErrEmptyToken
	}

	ctx := context.WithValue(r.Context(), gojwtnethttp.AuthorizationKey, token)
	return r.WithContext(ctx), nil
}

// GetCtxToken tries to get the raw token from the context
//
// Parameters:
//
//   - r: The HTTP request
//
// Returns:
//
//   - string: The raw token from the context
//   - error: An error if the token is not found or of an unexpected type
func GetCtxToken(r *http.Request) (string, error) {
	// Get the token from the context
	value := r.Context().Value(gojwtnethttp.AuthorizationKey)
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
