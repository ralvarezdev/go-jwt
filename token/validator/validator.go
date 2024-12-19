package validator

import (
	"github.com/golang-jwt/jwt/v5"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
)

// Validator does parsing and validation of JWT tokens
type (
	Validator interface {
		GetToken(rawToken string) (*jwt.Token, error)
		GetClaims(rawToken string) (*jwt.MapClaims, error)
		GetValidatedClaims(
			rawToken string,
			interception gojwtinterception.Interception,
		) (*jwt.MapClaims, error)
	}
)
