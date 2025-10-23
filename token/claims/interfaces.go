package claims

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// ClaimsValidator interface
	ClaimsValidator interface {
		ValidateClaims(
			ctx context.Context,
			claims jwt.MapClaims,
			token gojwttoken.Token,
		) (bool, error)
	}

	// TokenValidator interface
	TokenValidator interface {
		AddRefreshToken(
			ctx context.Context,
			id string,
			expiresAt time.Time,
		) error
		AddAccessToken(
			ctx context.Context,
			id string,
			parentRefreshTokenID string,
			expiresAt time.Time,
		) error
		RevokeToken(ctx context.Context, token gojwttoken.Token, id string) error
		IsTokenValid(ctx context.Context, token gojwttoken.Token, id string) (bool, error)
	}
)
