package grpc

import (
	"errors"
)

var (
	NilGRPCInterceptionsError             = errors.New("grpc interceptions map cannot be nil")
	MissingMetadataError                  = errors.New("missing metadata")
	AuthorizationMetadataInvalidError     = errors.New("authorization metadata invalid")
	AuthorizationMetadataNotProvidedError = errors.New("authorization metadata is not provided")
)
