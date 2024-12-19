package validator

import (
	"github.com/golang-jwt/jwt/v5"
	gojwtinterception "github.com/ralvarezdev/go-jwt/token/interception"
)

// ClaimsValidator interface
type ClaimsValidator interface {
	ValidateClaims(claims *jwt.MapClaims, interception gojwtinterception.Interception) (bool, error)
}
