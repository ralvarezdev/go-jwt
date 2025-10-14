package publisher

import (
	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
)

type (
	// Publisher is the interface for the JWT RabbitMQ publisher
	Publisher interface {
		SendTokenMessage(msg gojwtrabbitmq.TokensMessage) error
	}
)
