package validator

import (
	"errors"
)

var (
	ErrNilValidator            = errors.New("validator cannot be nil")
	ErrInvalidToken            = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidClaims           = errors.New("invalid claims")
	ErrNilClaims               = errors.New("claims cannot be nil")
	ErrNilClaimsValidator      = errors.New("claims validator cannot be nil")
)
