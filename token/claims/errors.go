package claims

import (
	"errors"
)

var (
	ErrIDClaimNotFound    = errors.New("id claim not found")
	ErrInvalidIDClaim     = errors.New("invalid id claim")
	ErrNilClaims          = errors.New("claims is nil")
	ErrNilTokenValidator  = errors.New("nil token validator")
	ErrNilClaimsValidator = errors.New("nil claims validator")
)
