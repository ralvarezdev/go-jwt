package sql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"

	godatabases "github.com/ralvarezdev/go-databases"
	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
	gojwtrabbitmqconsumer "github.com/ralvarezdev/go-jwt/rabbitmq/consumer"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// DefaultService is the default implementation of the Service interface
	DefaultService struct {
		handler  godatabasessql.Handler
		logger   *slog.Logger
		consumer gojwtrabbitmqconsumer.Consumer
		mutex    sync.Mutex
	}
)

// NewDefaultService creates a new DefaultService
//
// Parameters:
//
//   - handler: the SQL connection handler
//   - consumer: the RabbitMQ consumer
//   - logger: the logger (optional, can be nil)
//
// Returns:
//
//   - *DefaultService: the DefaultService instance
//   - error: an error if the data source or driver name is empty
func NewDefaultService(
	handler godatabasessql.Handler,
	consumer gojwtrabbitmqconsumer.Consumer,
	logger *slog.Logger,
) (*DefaultService, error) {
	// Check if the handler is nil
	if handler == nil {
		return nil, godatabases.ErrNilHandler
	}

	// Check if the consumer is nil
	if consumer == nil {
		return nil, gojwtrabbitmq.ErrNilConsumer
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "sqlite_jti_service"),
		)
	}

	return &DefaultService{
		handler:  handler,
		consumer: consumer,
		logger:   logger,
	}, nil
}

// Start starts the service to listen for messages and update the SQLite database
//
// Parameters:
//
//   - ctx: the context for managing cancellation and timeouts
//
// Returns:
//
//   - error: an error if the service could not be started
func (s *DefaultService) Start(ctx context.Context) error {
	// Check if the service is nil
	if s == nil {
		return sql.ErrConnDone
	}

	// Lock the mutex to ensure thread safety
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Connect to the database
	db, err := s.handler.Connect()
	if err != nil {
		if s.logger != nil {
			s.logger.Error(
				"Failed to connect to database",
				slog.String("error", err.Error()),
			)
		}
		return err
	}
	defer s.handler.Disconnect()

	// Ensure the table exists
	if _, err = db.Exec(CreateTableQuery); err != nil {
		return err
	}

	// Start the consumer
	tokensMessagesCh, err := s.consumer.ConsumeMessages(ctx)
	if err != nil {
		return err
	}

	// Listen for messages
	for {
		select {
		case <-ctx.Done():
			if s.logger != nil {
				s.logger.Info("Service context done, stopping service")
			}
			return nil
		case msg, ok := <-tokensMessagesCh:
			// Check if the channel is closed
			if !ok {
				if s.logger != nil {
					s.logger.Info("Message channel closed, stopping service")
				}
				return nil
			}

			// Process the message
			for _, issuedJTI := range msg.IssuedJTIs {
				if _, err = db.Exec(
					InsertTokenQuery,
					issuedJTI,
				); err != nil && s.logger != nil {
					s.logger.Error(
						"Failed to add JTI",
						slog.String("jti", issuedJTI),
						slog.String("error", err.Error()),
					)
				}
			}
			for _, revokedJTI := range msg.RevokedJTIs {
				if _, err = db.Exec(
					DeleteTokenQuery,
					revokedJTI,
				); err != nil && s.logger != nil {
					s.logger.Error(
						"Failed to remove JTI",
						slog.String("jti", revokedJTI),
						slog.String("error", err.Error()),
					)
				}
			}
		}
	}
}

// Validate checks if the given JTI exists in the database
//
// Parameters:
//
//   - jti: the JTI to validate
//
// Returns:
//
//   - bool: true if the JTI exists, false otherwise
//   - error: an error if the validation could not be performed
func (s *DefaultService) Validate(jti string) (bool, error) {
	// Check if the service is nil
	if s == nil {
		return false, sql.ErrConnDone
	}

	// Lock the mutex to ensure thread safety
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Connect to the database
	db, err := s.handler.Connect()
	if err != nil {
		if s.logger != nil {
			s.logger.Error(
				"Failed to connect to database",
				slog.String("error", err.Error()),
			)
		}
		return false, err
	}
	defer s.handler.Disconnect()

	// Check if the JTI exists
	var exists bool
	err = db.QueryRow(
		CheckTokenQuery,
		jti,
	).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if s.logger != nil {
			s.logger.Error(
				"Failed to validate JTI",
				slog.String("jti", jti),
				slog.String("error", err.Error()),
			)
		}
		return false, err
	}
	return exists, nil
}

// ValidateClaims validates the JWT claims
//
// Parameters:
//
//   - claims: the JWT claims to validate
//   - token: the token type
//
// Returns:
//
//   - bool: true if the claims are valid, false otherwise
//   - error: an error if the validation could not be performed
func (s *DefaultService) ValidateClaims(
	claims map[string]interface{},
	token gojwttoken.Token,
) (bool, error) {
	// Check if the service is nil
	if s == nil {
		return false, sql.ErrConnDone
	}

	// Extract the JTI from the claims
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		if s.logger != nil {
			s.logger.Error(
				"JTI claim is missing or invalid",
				slog.String("token", token.String()),
			)
		}
		return false, nil
	}

	// Validate the JTI
	return s.Validate(jti)
}
