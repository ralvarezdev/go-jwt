package rabbitmq

type (
	// TokensJTIMessage represents a message containing new and revoked JWT IDs
	TokensJTIMessage struct {
		NewJTIs     []string `json:"new_jtis"`
		RevokedJTIs []string `json:"revoked_jtis"`
	}
)
