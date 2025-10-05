package go_jwt

const (
	// BearerPrefix is the prefix for the bearer token
	BearerPrefix = "Bearer"
)

var (
	// CtxTokenClaimsKey is the key for the JWT context claims
	CtxTokenClaimsKey = "jwt_claims"

	// CtxTokenKey is the key for the JWT context token
	CtxTokenKey = "jwt"

	// IdClaim is the claim for the JWT ID
	IdClaim = "jti"

	// IsRefreshTokenClaim is the claim for refresh token
	IsRefreshTokenClaim = "irt"

	// SubjectClaim is the claim for the subject
	SubjectClaim = "sub"
)
