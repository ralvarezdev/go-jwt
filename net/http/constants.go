package http

type (
	// ContextKey is a type for context keys
	ContextKey string
)

var (
	// AuthorizationKey is the key of the authorization value in the context
	AuthorizationKey ContextKey = "authorization"
)
