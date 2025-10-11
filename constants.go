package go_jwt

const (
	// BearerPrefix is the prefix for the bearer token
	BearerPrefix = "Bearer"
)

var (
	// CtxRefreshTokenClaimsKey is the key for the refresh token context claims
	CtxRefreshTokenClaimsKey = "refresh_token_claims"

	// CtxAccessTokenClaimsKey is the key for the access token context claims
	CtxAccessTokenClaimsKey = "access_token_claims"

	// CtxRefreshTokenKey is the key for the refresh token context token
	CtxRefreshTokenKey = "refresh_token"

	// CtxAccessTokenKey is the key for the access token context token
	CtxAccessTokenKey = "access_token"

	// IDClaim is the claim for the JWT ID
	IDClaim = "jti"

	// IsRefreshTokenClaim is the claim for refresh token
	IsRefreshTokenClaim = "irt"

	// SubjectClaim is the claim for the subject
	SubjectClaim = "sub"
)
