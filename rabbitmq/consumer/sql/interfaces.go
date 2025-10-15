package sql

import (
	"context"

	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	gojwtclaims "github.com/ralvarezdev/go-jwt/token/claims"
)

type (
	// Service is the interface for the SQLite service for JWT IDs
	Service interface {
		gojwtclaims.Validator
		godatabasessql.Handler
		Start(ctx context.Context)
		IsRefreshTokenValid(jti string) (bool, error)
		IsAccessTokenValid(jti string) (bool, error)
	}
)
