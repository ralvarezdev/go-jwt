package rabbitmq

type (
	// TokensMessage represents a message containing new and revoked JWT IDs
	TokensMessage struct {
		IssuedJTIs  []string `json:"issued_jtis"`
		RevokedJTIs []string `json:"revoked_jtis"`
	}
)
