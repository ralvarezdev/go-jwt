package go_jwt

import "errors"

var (
	ErrUnableToParsePrivateKey            = errors.New("unable to parse private key")
	ErrUnableToParsePublicKey             = errors.New("unable to parse public key")
	ErrInvalidKeyType                     = errors.New("invalid key type")
	ErrMissingTokenInContext              = errors.New("no token in context")
	ErrMissingTokenClaimsInContext        = errors.New("no token claims in context")
	ErrUnexpectedTokenTypeInContext       = errors.New("unexpected token type in context")
	ErrUnexpectedTokenClaimsTypeInContext = errors.New("unexpected token claims type in context")
)
