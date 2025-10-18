package consumer

import (
	"context"

	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
)

type (
	// TokensMessagesConsumer is the interface for the JWT RabbitMQ tokens messages consumer
	TokensMessagesConsumer interface {
		GetChannel() <-chan *gojwtrabbitmq.TokensMessage
		ConsumeTokensMessages(ctx context.Context) error
	}

	// Consumer is the interface for the JWT RabbitMQ consumer
	Consumer interface {
		Open() error
		Close() error
		CreateTokensMessagesConsumer(ctx context.Context) (
			TokensMessagesConsumer,
			error,
		)
	}

	// Service is the interface for the SQLite service for JWT IDs
	Service interface {
		gojwttokenclaims.TokenValidator
		Start(ctx context.Context) error
	}
)
