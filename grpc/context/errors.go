package context

import (
	"errors"
)

var (
	ErrMissingToken              = errors.New("missing token")
	ErrUnexpectedTokenType       = errors.New("unexpected type")
	ErrMissingTokenClaims        = errors.New("missing token claims")
	ErrMissingTokenClaimsSubject = errors.New("missing token claims subject")
	ErrMissingTokenClaimsID      = errors.New("missing token claims id")
	ErrUnexpectedTokenClaimsType = errors.New("unexpected token claims type")
)
