package validator

import (
	"errors"
)

var (
	NilValidatorError            = errors.New("validator cannot be nil")
	InvalidTokenError            = errors.New("invalid token")
	UnexpectedSigningMethodError = errors.New("unexpected signing method")
	InvalidClaimsError           = errors.New("invalid claims")
	NilClaimsError               = errors.New("claims cannot be nil")
	NilClaimsValidatorError      = errors.New("claims validator cannot be nil")
)
