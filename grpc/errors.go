package grpc

import (
	"errors"
)

var (
	ErrNilGRPCInterceptions             = errors.New("grpc interceptions map cannot be nil")
	ErrMissingMetadata                  = errors.New("missing metadata")
	ErrAuthorizationMetadataInvalid     = errors.New("authorization metadata invalid")
	ErrAuthorizationMetadataNotProvided = errors.New("authorization metadata is not provided")
	ErrMissingToken              = errors.New("missing token")
	ErrUnexpectedTokenType       = errors.New("unexpected type")
	ErrMissingTokenClaims        = errors.New("missing token claims")
	ErrMissingTokenClaimsSubject = errors.New("missing token claims subject")
	ErrMissingTokenClaimsID      = errors.New("missing token claims id")
	ErrUnexpectedTokenClaimsType = errors.New("unexpected token claims type")
)
