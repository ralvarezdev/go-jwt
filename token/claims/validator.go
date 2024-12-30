package claims

import (
	"github.com/golang-jwt/jwt/v5"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
)

// Validator interface
type Validator interface {
	ValidateClaims(
		claims *jwt.MapClaims,
		interception gojwtinterception.Interception,
	) (bool, error)
}
