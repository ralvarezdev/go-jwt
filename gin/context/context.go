package context

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
)

// SetCtxTokenClaims sets the token claims in the context
//
// Parameters:
//
//   - ctx: The gin context
//   - key: The key to set the claims under
//   - claims: The token claims to set in the context
func SetCtxTokenClaims(ctx *gin.Context, key string, claims jwt.MapClaims) {
	ctx.Set(key, claims)
}

// SetCtxRefreshTokenClaims sets the refresh token claims in the context
//
// Parameters:
//
//   - ctx: The gin context
//   - claims: The refresh token claims to set in the context
func SetCtxRefreshTokenClaims(ctx *gin.Context, claims jwt.MapClaims) {
	ctx.Set(gojwt.CtxRefreshTokenClaimsKey, claims)
}

// SetCtxAccessTokenClaims sets the access token claims in the context
//
// Parameters:
//
//   - ctx: The gin context
//   - claims: The access token claims to set in the context
func SetCtxAccessTokenClaims(ctx *gin.Context, claims jwt.MapClaims) {
	ctx.Set(gojwt.CtxAccessTokenClaimsKey, claims)
}

// GetCtxTokenClaims tries to get the token claims from the context
//
// Parameters:
//
//   - ctx: The gin context
//   - key: The key to get the claims from
//
// Returns:
//
//   - jwt.MapClaims: The token claims from the context
//   - error: An error if the token claims are not found or of an unexpected type
func GetCtxTokenClaims(ctx *gin.Context, key string) (jwt.MapClaims, error) {
	// Get the token claims from the context
	value := ctx.Value(key)
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
//   - ctx: The gin context
//
// Returns:
//
//   - jwt.MapClaims: The refresh token claims from the context
//   - error: An error if the refresh token claims are not found or of an unexpected type
func GetCtxRefreshTokenClaims(ctx *gin.Context) (jwt.MapClaims, error) {
	return GetCtxTokenClaims(ctx, gojwt.CtxRefreshTokenClaimsKey)
}

// GetCtxAccessTokenClaims tries to get the access token claims from the context
//
// Parameters:
//
//   - ctx: The gin context
//
// Returns:
//
//   - jwt.MapClaims: The access token claims from the context
//   - error: An error if the access token claims are not found or of an unexpected type
func GetCtxAccessTokenClaims(ctx *gin.Context) (jwt.MapClaims, error) {
	return GetCtxTokenClaims(ctx, gojwt.CtxAccessTokenClaimsKey)
}

// SetCtxToken sets the raw token in the context
//
// Parameters:
//
//   - ctx: The gin context
//   - key: The key to set the token under
//   - token: The raw token to set in the context
func SetCtxToken(ctx *gin.Context, key, token string) {
	ctx.Set(key, token)
}

// SetCtxRefreshToken sets the raw refresh token in the context
//
// Parameters:
//
//   - ctx: The gin context
//   - token: The raw refresh token to set in the context
func SetCtxRefreshToken(ctx *gin.Context, token string) {
	SetCtxToken(ctx, gojwt.CtxRefreshTokenKey, token)
}

// SetCtxAccessToken sets the raw access token in the context
//
// Parameters:
//
//   - ctx: The gin context
//   - token: The raw access token to set in the context
func SetCtxAccessToken(ctx *gin.Context, token string) {
	SetCtxToken(ctx, gojwt.CtxAccessTokenKey, token)
}

// GetCtxToken tries to get the raw token from the context
//
// Parameters:
//
//   - ctx: The gin context
//   - key: The key to get the token from
//
// Returns:
//
//   - string: The raw token from the context
//   - error: An error if the token is not found or of an unexpected type
func GetCtxToken(ctx *gin.Context, key string) (string, error) {
	// Get the token from the context
	value := ctx.Value(key)
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
//   - ctx: The gin context
//
// Returns:
//
//   - string: The raw refresh token from the context
//   - error: An error if the refresh token is not found or of an unexpected type
func GetCtxRefreshToken(ctx *gin.Context) (string, error) {
	return GetCtxToken(ctx, gojwt.CtxRefreshTokenKey)
}

// GetCtxAccessToken tries to get the raw access token from the context
//
// Parameters:
//
//   - ctx: The gin context
//
// Returns:
//
//   - string: The raw access token from the context
//   - error: An error if the access token is not found or of an unexpected type
func GetCtxAccessToken(ctx *gin.Context) (string, error) {
	return GetCtxToken(ctx, gojwt.CtxAccessTokenKey)
}
