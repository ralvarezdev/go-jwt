package grpc

import (
	"errors"
)

var (
	ErrNilGRPCInterceptions             = errors.New("grpc interceptions map cannot be nil")
	ErrMissingMetadata                  = errors.New("missing metadata")
	ErrAuthorizationMetadataInvalid     = errors.New("authorization metadata invalid")
	ErrAuthorizationMetadataNotProvided = errors.New("authorization metadata is not provided")
)
