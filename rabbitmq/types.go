package rabbitmq

type (
	// TokensMessage represents a message containing new and revoked JWT IDs
	TokensMessage struct {
		IssuedRefreshTokensJTIs  []string `json:"issued_refresh_tokens_jtis"`
		RevokedRefreshTokensJTIs []string `json:"revoked_refresh_tokens_jtis"`
		IssuedAccessTokensJTIs   []string `json:"issued_access_tokens_jtis"`
		RevokedAccessTokensJTIs  []string `json:"revoked_access_tokens_jtis"`
	}
)
