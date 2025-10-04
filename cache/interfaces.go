package cache

import (
	"time"

	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
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
