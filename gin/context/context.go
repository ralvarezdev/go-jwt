package context

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtgin "github.com/ralvarezdev/go-jwt/gin"
)

// SetCtxRawToken sets the raw token in the context
func SetCtxRawToken(ctx *gin.Context, rawToken *string) {
	ctx.Set(gojwtgin.AuthorizationHeaderKey, *rawToken)
}

// SetCtxTokenClaims sets the token claims in the context
func SetCtxTokenClaims(ctx *gin.Context, claims *jwt.MapClaims) {
	ctx.Set(gojwt.CtxTokenClaimsKey, *claims)
}

// GetCtxRawToken tries to get the raw token from the context
func GetCtxRawToken(ctx *gin.Context) (string, error) {
	// Get the token from the context
	value := ctx.Value(gojwtgin.AuthorizationHeaderKey)
	if value == nil {
		return "", MissingTokenInContextError
	}

	// Check the type of the value
	rawToken, ok := value.(string)
	if !ok {
		return "", UnexpectedTokenTypeInContextError
	}

	return rawToken, nil
}

// GetCtxTokenClaims tries to get the token claims from the context
func GetCtxTokenClaims(ctx *gin.Context) (*jwt.MapClaims, error) {
	// Get the token claims from the context
	value := ctx.Value(gojwt.CtxTokenClaimsKey)
	if value == nil {
		return nil, MissingTokenClaimsInContextError
	}

	// Check the type of the value
	claims, ok := value.(jwt.MapClaims)
	if !ok {
		return nil, UnexpectedTokenClaimsTypeInContextError
	}

	return &claims, nil
}
