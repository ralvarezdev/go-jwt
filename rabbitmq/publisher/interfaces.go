package publisher

import (
	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
)

type (
	// Publisher is the interface for the JWT RabbitMQ publisher
	Publisher interface {
		Open() error
		Close() error
		PublishTokensMessage(msg gojwtrabbitmq.TokensMessage) error
	}
)
