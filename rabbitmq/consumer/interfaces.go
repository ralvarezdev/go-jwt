package consumer

import (
	"context"

	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
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
)
