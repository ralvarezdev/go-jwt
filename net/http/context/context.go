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
//   - claims: The token claims to set in the context
//
// Returns:
//
//   - *http.Request: The HTTP request with the token claims set in the context
func SetCtxTokenClaims(r *http.Request, claims jwt.MapClaims) *http.Request {
	ctx := context.WithValue(r.Context(), gojwt.CtxTokenClaimsKey, claims)
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
	value := r.Context().Value(gojwt.CtxTokenClaimsKey)
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
