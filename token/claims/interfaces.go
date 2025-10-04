package claims

import (
	"github.com/golang-jwt/jwt/v5"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// Validator interface
	Validator interface {
		ValidateClaims(
			claims jwt.MapClaims,
			token gojwttoken.Token,
		) (bool, error)
	}
)
