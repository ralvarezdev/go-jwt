package rabbitmq

import (
	"errors"
)

var (
	ErrNilConnection      = errors.New("nil rabbitmq connection")
	ErrEmptyQueueName     = errors.New("empty queue name")
	ErrNilChannel         = errors.New("nil rabbitmq channel")
	ErrNilPublisher       = errors.New("nil rabbitmq publisher")
	ErrNilConsumer        = errors.New("nil rabbitmq consumer")
	ErrNilMessage         = errors.New("nil rabbitmq message")
	ErrNilDeliveryChannel = errors.New("nil rabbitmq delivery channel")
)
