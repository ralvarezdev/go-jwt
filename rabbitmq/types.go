package rabbitmq

import (
	"time"
)

type (
	// TokenPair represents a pair of refresh and access token JTIs
	TokenPair struct {
		RefreshTokenID        string    `json:"refresh_token_id"`
		RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
		AccessTokenID         string    `json:"access_token_id"`
		AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	}

	// TokensMessage represents a message containing issued and revoked tokens
	TokensMessage struct {
		IssuedTokenPairs       []TokenPair `json:"issued_token_pairs"`
		RevokedRefreshTokensID []string    `json:"revoked_refresh_tokens_id"`
		RevokedAccessTokensID  []string    `json:"revoked_access_tokens_id"`
	}
)
