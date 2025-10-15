package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
)

type (
	// DefaultConsumer is the default implementation of the Consumer interface
	DefaultConsumer struct {
		conn                                    *amqp091.Connection
		ch                                      *amqp091.Channel
		queue                                   *amqp091.Queue
		queueName                               string
		logger                                  *slog.Logger
		period                                  time.Duration
		mutex                                   sync.Mutex
		tokensMessagesConsumerChannelBufferSize int
	}

	// DefaultTokensMessagesConsumer is the default implementation of the TokensMessagesConsumer interface
	DefaultTokensMessagesConsumer struct {
		period           time.Duration
		deliveryCh       <-chan amqp091.Delivery
		tokensMessagesCh chan gojwtrabbitmq.TokensMessage
		logger           *slog.Logger
	}
)

// NewDefaultConsumer creates a new DefaultConsumer
//
// Parameters:
//
//   - conn: the RabbitMQ connection
//   - queueName: the name of the queue
//   - period: the polling period
//   - tokensMessagesConsumerChannelBufferSize: the buffer size for the messages channel
//   - logger: the logger
//
// Returns:
//
//   - *DefaultConsumer: the DefaultConsumer instance
//   - error: an error if the connection is nil
func NewDefaultConsumer(
	conn *amqp091.Connection,
	queueName string,
	period time.Duration,
	tokensMessagesConsumerChannelBufferSize int,
	logger *slog.Logger,
) (*DefaultConsumer, error) {
	// Check if the connection is nil
	if conn == nil {
		return nil, gojwtrabbitmq.ErrNilConnection
	}

	// Check if the queue name is empty
	if queueName == "" {
		return nil, gojwtrabbitmq.ErrEmptyQueueName
	}

	// Check if the period is valid
	if period <= 0 {
		period = DefaultTickerInterval
	}

	// Check if the consumerMessagesBufferSize is valid
	if tokensMessagesConsumerChannelBufferSize <= 0 {
		tokensMessagesConsumerChannelBufferSize = DefaultTokensMessageConsumerChannelBufferSize
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "jwt_rabbitmq_consumer"),
		)
	}

	// Create a new consumer instance
	consumer := &DefaultConsumer{
		conn:                                    conn,
		logger:                                  logger,
		queueName:                               queueName,
		period:                                  period,
		tokensMessagesConsumerChannelBufferSize: tokensMessagesConsumerChannelBufferSize,
	}
	return consumer, nil
}

// NewDefaultTokensMessagesConsumer creates a new DefaultTokensMessagesConsumer
//
// Parameters:
//
//   - deliveryCh: the RabbitMQ delivery channel
//   - bufferSize: the buffer size for the messages channel
//   - period: the polling period
//   - logger: the logger
//
// Returns:
//
//   - *DefaultTokensMessagesConsumer: the DefaultTokensMessagesConsumer instance
//   - error: an error if the delivery channel is nil
func NewDefaultTokensMessagesConsumer(
	deliveryCh <-chan amqp091.Delivery,
	bufferSize int,
	period time.Duration,
	logger *slog.Logger,
) (*DefaultTokensMessagesConsumer, error) {
	// Check if the delivery channel is nil
	if deliveryCh == nil {
		return nil, gojwtrabbitmq.ErrNilDeliveryChannel
	}

	// Check if the buffer size is valid
	if bufferSize <= 0 {
		bufferSize = DefaultTokensMessageConsumerChannelBufferSize
	}

	// Check if the period is valid
	if period <= 0 {
		period = DefaultTickerInterval
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "jwt_rabbitmq_tokens_messages_consumer"),
		)
	}

	return &DefaultTokensMessagesConsumer{
		deliveryCh:       deliveryCh,
		tokensMessagesCh: make(chan gojwtrabbitmq.TokensMessage, bufferSize),
		period:           period,
		logger:           logger,
	}, nil
}

// Open opens a RabbitMQ channel
//
// Returns:
//
//   - error: an error if the channel could not be opened
func (d *DefaultConsumer) Open() error {
	// Check if the consumer is nil
	if d == nil {
		return gojwtrabbitmq.ErrNilConsumer
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if the consumer is already open
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
	d.logger.Info("Consumer channel opened")
	return nil
}

// Close closes the RabbitMQ channel and connection
//
// Returns:
//
//   - error: an error if the channel could not be closed
func (d *DefaultConsumer) Close() error {
	// Check if the consumer is nil
	if d == nil {
		return gojwtrabbitmq.ErrNilConsumer
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
	d.logger.Info("Consumer channel closed")
	return nil
}

// CreateTokensMessagesConsumer creates a tokens messages consumer
//
// Parameters:
//
//   - ctx: the context
//
// Returns:
//
//   - TokensMessagesConsumer: the tokens messages consumer
//   - error: an error if the consumer could not be created
func (d *DefaultConsumer) CreateTokensMessagesConsumer(ctx context.Context) (
	TokensMessagesConsumer,
	error,
) {
	if d == nil {
		return nil, gojwtrabbitmq.ErrNilConsumer
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Ensure the channel is open
	if d.ch == nil {
		if err := d.Open(); err != nil {
			return nil, err
		}
	}

	// Create a channel to receive messages
	deliveryCh, err := gojwtrabbitmq.CreateConsumeTokensMessageDeliveryChWithCtx(
		ctx,
		d.ch,
		d.queue.Name,
	)
	if err != nil {
		return nil, err
	}

	// Create the tokens messages consumer
	consumer, err := NewDefaultTokensMessagesConsumer(
		deliveryCh,
		d.tokensMessagesConsumerChannelBufferSize,
		d.period,
		d.logger,
	)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// GetChannel returns the tokens message consumer channel
//
// Returns:
//
//   - <-chan gojwtrabbitmq.TokensMessage: the tokens message consumer channel
func (d DefaultTokensMessagesConsumer) GetChannel() <-chan gojwtrabbitmq.TokensMessage {
	return d.tokensMessagesCh
}

// ConsumeTokensMessages consumes a tokens message and sends it to the channel
//
// Parameters:
//
//   - ctx: the context
//
// Returns:
func (d DefaultTokensMessagesConsumer) ConsumeTokensMessages(
	ctx context.Context,
) error {
	// Create the ticker to poll the queue periodically
	ticker := time.NewTicker(d.period)

	for {
		select {
		case <-ctx.Done():
			if d.logger != nil {
				d.logger.Info("Context done. Exiting consume loop.")
			}
			return ctx.Err()
		case <-ticker.C:
			// Poll the queue for messages
			for msg := range d.deliveryCh {
				var parsedMsg gojwtrabbitmq.TokensMessage
				if err := json.Unmarshal(msg.Body, &parsedMsg); err != nil {
					if d.logger != nil {
						d.logger.Error(
							"Error decoding message",
							slog.String("error", err.Error()),
						)
					}
					continue
				}

				// Send the message to the channel
				d.tokensMessagesCh <- parsedMsg
			}
		}
	}
}
