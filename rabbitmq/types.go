package rabbitmq

type (
	// TokenPair represents a pair of refresh and access token JTIs
	TokenPair struct {
		RefreshTokenJTI string `json:"refresh_token_jti"`
		AccessTokenJTI  string `json:"access_token_jti"`
	}

	// TokensMessage represents a message containing issued and revoked tokens
	TokensMessage struct {
		IssuedTokenPairs         []TokenPair `json:"issued_token_pairs"`
		RevokedRefreshTokensJTIs []string    `json:"revoked_refresh_tokens_jtis"`
		RevokedAccessTokensJTIs  []string    `json:"revoked_access_tokens_jtis"`
	}
)
