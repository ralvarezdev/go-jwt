package grpc

type (
	// ContextKey is a type for context keys
	ContextKey string
)

var (
	// AuthorizationTokenIdx is the index for the authorization token
	AuthorizationTokenIdx = 0

	// AuthorizationKey is the key for the authorization token in the context
	AuthorizationKey ContextKey = "authorization"
)
