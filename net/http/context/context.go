package context

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtgin "github.com/ralvarezdev/go-jwt/gin"
	gojwthttp "github.com/ralvarezdev/go-jwt/net/http"
	"net/http"
)

// SetCtxRawToken sets the raw token in the context
func SetCtxRawToken(r *http.Request, rawToken *string) {
	ctx := context.WithValue(
		r.Context(),
		gojwthttp.AuthorizationHeaderKey,
		*rawToken,
	)
	r = r.WithContext(ctx)
}

// SetCtxTokenClaims sets the token claims in the context
func SetCtxTokenClaims(r *http.Request, claims *jwt.MapClaims) {
	ctx := context.WithValue(r.Context(), gojwt.CtxTokenClaimsKey, *claims)
	r = r.WithContext(ctx)
}

// GetCtxRawToken tries to get the raw token from the context
func GetCtxRawToken(r *http.Request) (string, error) {
	// Get the token from the context
	value := r.Context().Value(gojwtgin.AuthorizationHeaderKey)
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
