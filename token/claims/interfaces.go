package claims

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// ClaimsValidator interface
	ClaimsValidator interface {
		ValidateClaims(
			claims jwt.MapClaims,
			token gojwttoken.Token,
		) (bool, error)
	}

	// TokenValidator interface
	TokenValidator interface {
		AddRefreshToken(
			id string,
			expiresAt time.Time,
		) error
		AddAccessToken(
			id string,
			parentRefreshTokenID string,
			expiresAt time.Time,
		) error
		RevokeToken(token gojwttoken.Token, id string) error
		IsTokenValid(token gojwttoken.Token, id string) (bool, error)
	}
)
