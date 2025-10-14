package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rabbitmq/amqp091-go"
	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
)

type (
	// DefaultConsumer is the default implementation of the Consumer interface
	DefaultConsumer struct {
		conn                       *amqp091.Connection
		ch                         *amqp091.Channel
		queue                      *amqp091.Queue
		queueName                  string
		logger                     *slog.Logger
		period                     time.Duration
		mutex                      sync.Mutex
		isTickerRunning            atomic.Bool
		tickerStopCh               chan struct{}
		consumerMessagesBufferSize int
	}
)

// NewDefaultConsumer creates a new DefaultConsumer
//
// Parameters:
//
//   - conn: the RabbitMQ connection
//   - queueName: the name of the queue
//   - period: the polling period
//   - consumerMessagesChannelBufferSize: the size of the messages buffer channel
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
	consumerMessagesChannelBufferSize int,
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
	if consumerMessagesChannelBufferSize <= 0 {
		consumerMessagesChannelBufferSize = DefaultConsumerMessagesChannelBufferSize
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "jwt_rabbitmq_consumer"),
		)
	}

	// Create a new consumer instance
	consumer := &DefaultConsumer{
		conn:                       conn,
		logger:                     logger,
		queueName:                  queueName,
		period:                     period,
		consumerMessagesBufferSize: consumerMessagesChannelBufferSize,
	}
	return consumer, nil
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

	// Check if the ticker is already running
	if d.isTickerRunning.Load() {
		if d.logger != nil {
			d.logger.Warn("Ticker is already running. Skipping open operation.")
		}
		return nil
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
	q, err := gojwtrabbitmq.DeclareJTIQueue(d.ch, d.queueName)
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

	// Create a channel to signal the ticker to stop
	d.tickerStopCh = make(chan struct{})
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

	// Check if the ticker is running and stop it
	if d.isTickerRunning.Load() {
		d.isTickerRunning.Store(false)

		// Signal the ticker to stop
		close(d.tickerStopCh)

		if d.logger != nil {
			d.logger.Info("Ticker stopped.")
		}
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

// ConsumeMessages consumes messages from the RabbitMQ queue
//
// Parameters:
//
//   - ctx: the context
//
// Returns:
//
//   - <-chan gojwtrabbitmq.TokensMessage: a channel to receive the messages
//   - error: an error if the consumer is nil or the channel is not open
func (d *DefaultConsumer) ConsumeMessages(ctx context.Context) (
	<-chan gojwtrabbitmq.TokensMessage,
	error,
) {
	if d == nil {
		return nil, gojwtrabbitmq.ErrNilConsumer
	}

	// Check if the ticker is already running
	if d.isTickerRunning.Load() {
		if d.logger != nil {
			d.logger.Warn("Ticker is already running. Skipping consume operation.")
		}
		return nil, nil
	}

	// Mark the ticker as running
	d.isTickerRunning.Store(true)

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Ensure the channel is open
	if d.ch == nil {
		if err := d.Open(); err != nil {
			return nil, err
		}
	}

	// Create the channel to return messages
	msgCh := make(
		chan gojwtrabbitmq.TokensMessage,
		d.consumerMessagesBufferSize,
	)

	// Create the ticker to poll the queue periodically
	ticker := time.NewTicker(d.period)

	// Declare the JTI queue
	q, err := gojwtrabbitmq.DeclareJTIQueue(d.ch, d.queueName)
	if err != nil {
		return nil, err
	}

	// Create a channel to receive messages
	deliveryCh, err := gojwtrabbitmq.CreateConsumeJTIDeliveryChWithCtx(
		ctx,
		d.ch,
		q.Name,
	)
	if err != nil {
		return nil, err
	}

	for {
		select {
		case <-ctx.Done():
			if d.logger != nil {
				d.logger.Info("Context done. Exiting consume loop.")
			}
			return nil, ctx.Err()
		case <-d.tickerStopCh:
			if d.logger != nil {
				d.logger.Info("Ticker stop signal received. Exiting consume loop.")
			}
			return nil, nil
		case <-ticker.C:
			// Poll the queue for messages
			for msg := range deliveryCh {
				var parsedMsg gojwtrabbitmq.TokensMessage
				if err = json.Unmarshal(msg.Body, &parsedMsg); err != nil {
					if d.logger != nil {
						d.logger.Error(
							"Error decoding message",
							slog.String("error", err.Error()),
						)
					}
					continue
				}

				// Send the message to the channel
				msgCh <- parsedMsg
			}
		}
	}
}
