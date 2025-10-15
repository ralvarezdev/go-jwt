package publisher

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
)

type (
	// DefaultPublisher is the default implementation of the Publisher interface
	DefaultPublisher struct {
		conn      *amqp091.Connection
		ch        *amqp091.Channel
		queue     *amqp091.Queue
		queueName string
		logger    *slog.Logger
		mutex     sync.Mutex
	}
)

// NewDefaultPublisher creates a new DefaultPublisher
//
// Parameters:
//
//   - conn: the RabbitMQ connection
//   - queueName: the name of the queue
//   - logger: the logger
//
// Returns:
//
//   - *DefaultPublisher: the DefaultPublisher instance
//   - error: an error if the connection is nil
func NewDefaultPublisher(
	conn *amqp091.Connection,
	queueName string,
	logger *slog.Logger,
) (*DefaultPublisher, error) {
	// Check if the connection is nil
	if conn == nil {
		return nil, gojwtrabbitmq.ErrNilConnection
	}

	// Check if the queue name is empty
	if queueName == "" {
		return nil, gojwtrabbitmq.ErrEmptyQueueName
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "jwt_rabbitmq_publisher"),
		)
	}

	// Create a new publisher instance
	publisher := &DefaultPublisher{
		conn:      conn,
		queueName: queueName,
		logger:    logger,
	}
	return publisher, nil
}

// Open opens a RabbitMQ channel
//
// Returns:
//
//   - error: an error if the channel could not be opened
func (d *DefaultPublisher) Open() error {
	// Check if the publisher is nil
	if d == nil {
		return gojwtrabbitmq.ErrNilPublisher
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if the publisher is already open
	if d.ch != nil {
		return nil
	}

	// Create the channel
	ch, err := d.conn.Channel()
	if err != nil {
		d.logger.Error(
			"Failed to open a channel",
			slog.String("error", err.Error()),
		)
		return err
	}

	// Set the channel
	d.ch = ch

	// Declare the queue
	q, err := gojwtrabbitmq.DeclareTokensMessageQueue(d.ch, d.queueName)
	if err != nil {
		d.logger.Error(
			"Failed to declare a queue",
			slog.String("queue_name", d.queueName),
			slog.String("error", err.Error()),
		)
		return err
	}
	d.queue = q
	d.logger.Info("Publisher channel opened")
	return nil
}

// Close closes the RabbitMQ channel and connection
//
// Returns:
//
//   - error: an error if the channel could not be closed
func (d *DefaultPublisher) Close() error {
	// Check if the publisher is nil
	if d == nil {
		return gojwtrabbitmq.ErrNilPublisher
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.ch == nil {
		return nil
	}

	// Close the channel
	if err := d.ch.Close(); err != nil {
		d.logger.Error(
			"Failed to close channel", slog.String("error", err.Error()),
		)
		return err
	}

	// Set the channel to nil
	d.ch = nil
	d.logger.Info("Publisher channel closed")
	return nil
}

// PublishTokensMessage publishes a tokens message to the RabbitMQ queue
//
// Parameters:
//
//   - msg: the tokens message to publish
//
// Returns:
//
//   - error: an error if the message could not be published
func (d *DefaultPublisher) PublishTokensMessage(msg *gojwtrabbitmq.TokensMessage) error {
	// Check if the publisher is nil
	if d == nil {
		return gojwtrabbitmq.ErrNilPublisher
	}

	// Check if the message is nil
	if msg == nil {
		return gojwtrabbitmq.ErrNilMessage
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()

	// Ensure the channel is open
	if d.ch == nil {
		if err := d.Open(); err != nil {
			d.mutex.Unlock()
			return err
		}
	}
	d.mutex.Unlock()

	// Marshal the message to JSON
	body, err := json.Marshal(msg)
	if err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to marshal message",
				slog.String("error", err.Error()),
			)
		}
		return err
	}

	// Publish the message
	err = d.ch.Publish(
		"",
		d.queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to publish message",
				slog.String("error", err.Error()),
			)
		}
		return err
	}
	if d.logger != nil {
		d.logger.Debug("Message published", slog.String("body", string(body)))
	}
	return nil
}
