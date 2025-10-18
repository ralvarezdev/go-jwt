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
		Set(
			token gojwttoken.Token,
			id string,
			isValid bool,
			expiresAt time.Time,
		) error
		Revoke(token gojwttoken.Token, id string) error
		IsValid(token gojwttoken.Token, id string) (bool, error)
	}
)
