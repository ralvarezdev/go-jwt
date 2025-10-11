package context

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtgin "github.com/ralvarezdev/go-jwt/gin"
)

// SetCtxTokenClaims sets the token claims in the context
//
// Parameters:
//
//   - ctx: The gin context
//   - claims: The token claims to set in the context
func SetCtxTokenClaims(ctx *gin.Context, claims jwt.MapClaims) {
	ctx.Set(gojwtgin.AuthorizationKey, claims)
}

// GetCtxTokenClaims tries to get the token claims from the context
//
// Parameters:
//
//   - ctx: The gin context
//
// Returns:
//
//   - jwt.MapClaims: The token claims from the context
//   - error: An error if the token claims are not found or of an unexpected type
func GetCtxTokenClaims(ctx *gin.Context) (jwt.MapClaims, error) {
	// Get the token claims from the context
	value := ctx.Value(gojwtgin.AuthorizationKey)
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
//   - ctx: The gin context
//   - token: The raw token to set in the context
func SetCtxToken(ctx *gin.Context, token string) {
	ctx.Set(gojwtgin.AuthorizationKey, token)
}

// GetCtxToken tries to get the raw token from the context
//
// Parameters:
//
//   - ctx: The gin context
//
// Returns:
//
//   - string: The raw token from the context
//   - error: An error if the token is not found or of an unexpected type
func GetCtxToken(ctx *gin.Context) (string, error) {
	// Get the token from the context
	value := ctx.Value(gojwtgin.AuthorizationKey)
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
