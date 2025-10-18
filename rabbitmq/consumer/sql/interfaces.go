package sql

import (
	"context"

	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
	gojwtclaims "github.com/ralvarezdev/go-jwt/token/claims"
)

type (
	// Service is the interface for the SQLite service for JWT IDs
	Service interface {
		gojwtclaims.Validator
		godatabasessql.Handler
		Start(ctx context.Context) error
		InsertRefreshTokens(jtis ...string) error
		InsertAccessTokens(tokenPairs ...gojwtrabbitmq.TokenPair) error
		RevokeRefreshTokens(jtis ...string) error
		RevokeAccessTokens(jtis ...string) error
		RevokeAccessTokensByRefreshTokens(jtis ...string) error
		IsRefreshTokenValid(jti string) (bool, error)
		IsAccessTokenValid(jti string) (bool, error)
	}
)
