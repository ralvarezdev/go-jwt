package gojwt

type (
	// CtxKey is the type for the context keys
	CtxKey string
)

const (
	// BearerPrefix is the prefix for the bearer token
	BearerPrefix = "Bearer"
)

var (
	// CtxTokenClaimsKey is the key for the token claims to be set to the context
	CtxTokenClaimsKey CtxKey = "token_claims"

	// CtxTokenKey is the key for the token to be set to the context
	CtxTokenKey CtxKey = "token"

	// IDClaim is the claim for the JWT ID
	IDClaim = "jti"

	// IsRefreshTokenClaim is the claim for refresh token
	IsRefreshTokenClaim = "irt"

	// SubjectClaim is the claim for the subject
	SubjectClaim = "sub"
)
