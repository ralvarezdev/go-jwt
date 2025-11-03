package gojwt

const (
	// BearerPrefix is the prefix for the bearer token
	BearerPrefix = "Bearer"
)

var (
	// CtxTokenClaimsKey is the key for the token claims to be set to the context
	CtxTokenClaimsKey = "token_claims"

	// CtxTokenKey is the key for the token to be set to the context
	CtxTokenKey = "token"

	// IDClaim is the claim for the JWT ID
	IDClaim = "jti"

	// IsRefreshTokenClaim is the claim for refresh token
	IsRefreshTokenClaim = "irt"

	// SubjectClaim is the claim for the subject
	SubjectClaim = "sub"
)
