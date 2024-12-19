package metadata

import (
	"context"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtgrpc "github.com/ralvarezdev/go-jwt/grpc"
	goloadergcloud "github.com/ralvarezdev/go-loader/cloud/gcloud"
	"google.golang.org/grpc/metadata"
	"strings"
)

type (
	// Field is a field in the metadata
	Field struct {
		Key   string
		Value string
	}

	// CtxMetadata is the metadata for the context
	CtxMetadata struct {
		Fields []Field
	}
)

// NewCtxMetadata creates a new CtxMetadata
func NewCtxMetadata(metadataFields *map[string]string) (*CtxMetadata, error) {
	// Check if the metadata fields are nil
	if metadataFields == nil {
		return nil, NilMetadataFieldsError
	}

	// Add the metadata fields
	var fields []Field
	for key, value := range *metadataFields {
		fields = append(
			fields,
			Field{Key: strings.ToLower(key), Value: value},
		)
	}

	return &CtxMetadata{
		Fields: fields,
	}, nil
}

// NewUnauthenticatedCtxMetadata creates a new unauthenticated CtxMetadata
func NewUnauthenticatedCtxMetadata(gcloudToken string) (*CtxMetadata, error) {
	return NewCtxMetadata(
		&map[string]string{
			goloadergcloud.AuthorizationMetadataKey: gojwt.BearerPrefix + " " + gcloudToken,
		},
	)
}

// NewAuthenticatedCtxMetadata creates a new authenticated CtxMetadata
func NewAuthenticatedCtxMetadata(
	gcloudToken string, jwtToken string,
) (*CtxMetadata, error) {
	return NewCtxMetadata(
		&map[string]string{
			goloadergcloud.AuthorizationMetadataKey: gojwt.BearerPrefix + " " + gcloudToken,
			gojwtgrpc.AuthorizationMetadataKey:      gojwt.BearerPrefix + " " + jwtToken,
		},
	)
}

// GetCtxWithMetadata gets the context with the metadata
func GetCtxWithMetadata(
	ctxMetadata *CtxMetadata, ctx context.Context,
) context.Context {
	// Check if the context metadata is nil
	if ctxMetadata == nil {
		return ctx
	}

	// Create metadata
	md := metadata.Pairs()

	// Add the metadata to the context
	for _, field := range ctxMetadata.Fields {
		md.Append(field.Key, field.Value)
	}
	return metadata.NewOutgoingContext(ctx, md)
}

// AppendGCloudTokenToOutgoingContext appends the GCloud token to the outgoing context
func AppendGCloudTokenToOutgoingContext(
	ctx context.Context, gcloudToken string,
) context.Context {
	return metadata.AppendToOutgoingContext(
		ctx,
		goloadergcloud.AuthorizationMetadataKey,
		gojwt.BearerPrefix+" "+gcloudToken,
	)
}
