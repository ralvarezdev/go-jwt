package context

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
	"net/http"
)

// SetCtxTokenClaims sets the token claims in the context
func SetCtxTokenClaims(r *http.Request, claims *jwt.MapClaims) *http.Request {
	ctx := context.WithValue(r.Context(), gojwt.CtxTokenClaimsKey, *claims)
	return r.WithContext(ctx)
}

// GetCtxTokenClaims tries to get the token claims from the context
func GetCtxTokenClaims(r *http.Request) (*jwt.MapClaims, error) {
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

	return &claims, nil
}
