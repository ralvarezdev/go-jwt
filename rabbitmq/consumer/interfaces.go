package consumer

import (
	"context"

	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
)

type (
	// Consumer is the interface for the JWT RabbitMQ consumer
	Consumer interface {
		Open() error
		Close() error
		ConsumeMessages(ctx context.Context) (
			<-chan gojwtrabbitmq.TokensMessage,
			error,
		)
	}
)
