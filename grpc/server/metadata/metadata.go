package metadata

import (
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtgrpc "github.com/ralvarezdev/go-jwt/grpc"
	goloadergcloud "github.com/ralvarezdev/go-loader/cloud/gcloud"
	"google.golang.org/grpc/metadata"
	"strings"
)

// GetTokenFromMetadata gets the token from the metadata
func GetTokenFromMetadata(md metadata.MD, tokenKey string) (string, error) {
	// Get the authorization from the metadata
	authorization := md.Get(tokenKey)
	tokenIdx := gojwtgrpc.TokenIdx.Int()
	if len(authorization) <= tokenIdx {
		return "", gojwtgrpc.AuthorizationMetadataNotProvidedError
	}

	// Get the authorization value from the metadata
	authorizationValue := authorization[tokenIdx]

	// Split the authorization value by space
	authorizationFields := strings.Split(authorizationValue, " ")

	// Check if the authorization value is valid
	if len(authorizationFields) != 2 || authorizationFields[0] != gojwt.BearerPrefix {
		return "", gojwtgrpc.AuthorizationMetadataInvalidError
	}

	return authorizationFields[1], nil
}

// GetAuthorizationTokenFromMetadata gets the authorization token from the metadata
func GetAuthorizationTokenFromMetadata(md metadata.MD) (string, error) {
	return GetTokenFromMetadata(md, gojwtgrpc.AuthorizationMetadataKey)
}

// GetGCloudAuthorizationTokenFromMetadata gets the GCloud authorization token from the metadata
func GetGCloudAuthorizationTokenFromMetadata(md metadata.MD) (string, error) {
	return GetTokenFromMetadata(md, goloadergcloud.AuthorizationMetadataKey)
}
