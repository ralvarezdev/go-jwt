package sql

import (
	"context"

	gojwtclaims "github.com/ralvarezdev/go-jwt/token/claims"
)

type (
	// Service is the interface for the SQLite service for JWT IDs
	Service interface {
		gojwtclaims.Validator
		Start(ctx context.Context)
		Validate(jti string) (bool, error)
	}
)
