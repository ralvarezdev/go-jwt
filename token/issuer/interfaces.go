package issuer

import (
	"github.com/golang-jwt/jwt/v5"
)

type (
	// Issuer is the interface for JWT tokens issuing
	Issuer interface {
		IssueToken(claims jwt.Claims) (string, error)
	}
)
