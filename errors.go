package go_jwt

import "errors"

var (
	UnableToParsePrivateKeyError = errors.New("unable to parse private key")
	UnableToParsePublicKeyError  = errors.New("unable to parse public key")
	InvalidKeyTypeError          = errors.New("invalid key type")
)
