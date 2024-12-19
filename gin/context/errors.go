package context

import "errors"

var (
	MissingTokenInContextError              = errors.New("no token in context")
	MissingTokenClaimsInContextError        = errors.New("no token claims in context")
	UnexpectedTokenTypeInContextError       = errors.New("unexpected token type in context")
	UnexpectedTokenClaimsTypeInContextError = errors.New("unexpected token claims type in context")
)
