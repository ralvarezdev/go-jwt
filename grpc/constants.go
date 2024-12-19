package grpc

type AuthorizationIdx int

const (
	TokenIdx AuthorizationIdx = iota
)

// Int returns the integer value of the index
func (i AuthorizationIdx) Int() int {
	return int(i)
}

const (
	// AuthorizationMetadataKey is the key for the authorization metadata
	AuthorizationMetadataKey = "authorization"
)
