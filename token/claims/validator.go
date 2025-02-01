package claims

import (
	"github.com/golang-jwt/jwt/v5"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

// Validator interface
type Validator interface {
	ValidateClaims(
		claims *jwt.MapClaims,
		token gojwttoken.Token,
	) (bool, error)
}
