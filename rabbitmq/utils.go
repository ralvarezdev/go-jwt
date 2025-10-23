package rabbitmq

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

// DeclareTokensMessageQueue creates a durable queue for handling tokens messages
//
// Parameters:
//
//   - ch: the RabbitMQ channel
//   - queueName: the name of the queue
//
// Returns:
//
//   - *amqp091.Queue: the declared queue
//   - error: an error if the queue could not be declared
func DeclareTokensMessageQueue(ch *amqp091.Channel, queueName string) (
	*amqp091.Queue,
	error,
) {
	// Check if the channel is nil
	if ch == nil {
		return nil, ErrNilChannel
	}

	// Check if the queue name is empty
	if queueName == "" {
		return nil, ErrEmptyQueueName
	}

	// Declare a durable queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// CreateConsumeTokensMessageDeliveryChWithCtx sets up a consumer to receive messages from the specified queue using the
// provided context
//
// Parameters:
//
//   - ctx: the context for managing cancellation and timeouts
//   - ch: the RabbitMQ channel
//   - queueName: the name of the queue
//
// Returns:
//
//   - <-chan amqp091.Delivery: a channel to receive messages
//   - error: an error if the consumer could not be set up
func CreateConsumeTokensMessageDeliveryChWithCtx(
	ctx context.Context,
	ch *amqp091.Channel,
	queueName string,
) (<-chan amqp091.Delivery, error) {
	// Check if the channel is nil
	if ch == nil {
		return nil, ErrNilChannel
	}

	// Check if the queue name is empty
	if queueName == "" {
		return nil, ErrEmptyQueueName
	}

	// Declare the queue to ensure it exists
	deliveryCh, err := ch.ConsumeWithContext(
		ctx,
		"",
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	return deliveryCh, nil
}
