package context

import "errors"

var (
	MissingTokenError              = errors.New("missing token")
	UnexpectedTokenTypeError       = errors.New("unexpected type")
	MissingTokenClaimsError        = errors.New("missing token claims")
	MissingTokenClaimsSubjectError = errors.New("missing token claims subject")
	MissingTokenClaimsIdError      = errors.New("missing token claims id")
	UnexpectedTokenClaimsTypeError = errors.New("unexpected token claims type")
)
