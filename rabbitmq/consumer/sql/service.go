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
	"golang.org/x/sync/errgroup"
)

type (
	// DefaultService is the default implementation of the Service interface
	DefaultService struct {
		godatabasessql.Handler
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
		Handler:  handler,
		consumer: consumer,
		logger:   logger,
	}, nil
}

// Connect opens the database connection
//
// Returns:
//
//   - error: an error if the connection could not be opened
func (d *DefaultService) Connect() error {
	// Check if the service is nil
	if d == nil {
		return godatabases.ErrNilHandler
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Connect to the database
	db, err := d.Handler.Connect()
	if err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to connect to database",
				slog.String("error", err.Error()),
			)
		}
		return err
	}

	// Ensure the tables exist
	if _, err = db.Exec(CreateRefreshTokensTableQuery); err != nil {
		return err
	}
	if _, err = db.Exec(CreateAccessTokensTableQuery); err != nil {
		return err
	}

	return nil
}

// Disconnect closes the database connection
//
// Returns:
//
//   - error: an error if the connection could not be closed
func (d *DefaultService) Disconnect() error {
	// Check if the service is nil
	if d == nil {
		return godatabases.ErrNilHandler
	}

	// Lock the mutex to ensure thread safety
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Disconnect from the database
	if err := d.Handler.Disconnect(); err != nil {
		if d.logger != nil {
			d.logger.Error(
				"Failed to disconnect from database",
				slog.String("error", err.Error()),
			)
		}
		return err
	}

	return nil
}

// DB is a helper function to get the database connection
//
// Returns:
//
//   - *sql.DB: the database connection
func (d *DefaultService) DB() (*sql.DB, error) {
	// Lock the mutex to ensure thread safety
	d.mutex.Lock()

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		d.mutex.Unlock()
		if d.logger != nil {
			d.logger.Error(
				"Failed to get database connection",
				slog.String("error", err.Error()),
			)
		}
		return nil, err
	}
	d.mutex.Unlock()

	return db, nil
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
func (d *DefaultService) Start(ctx context.Context) error {
	// Check if the service is nil
	if d == nil {
		return sql.ErrConnDone
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return err
	}

	// Start the consumer
	tokensMessagesConsumer, err := d.consumer.CreateTokensMessagesConsumer(ctx)
	if err != nil {
		return err
	}

	// Create an error group to handle errors from the consumer
	eg := errgroup.Group{}

	// Start the consumer in a separate goroutine
	eg.Go(
		func() error {
			return tokensMessagesConsumer.ConsumeTokensMessages(ctx)
		},
	)

	// Listen for messages
	eg.Go(
		func() error {
			// Get the messages channel
			tokensMessagesCh := tokensMessagesConsumer.GetChannel()

			for {
				select {
				case <-ctx.Done():
					if d.logger != nil {
						d.logger.Info("Service context done, stopping service")
					}
					return nil
				case msg, ok := <-tokensMessagesCh:
					// Check if the channel is closed
					if !ok {
						if d.logger != nil {
							d.logger.Info("Message channel closed, stopping service")
						}
						return nil
					}

					// Process the message
					for _, issuedRefreshTokenJTI := range msg.IssuedRefreshTokensJTIs {
						if _, err = db.Exec(
							InsertRefreshTokenQuery,
							issuedRefreshTokenJTI,
						); err != nil && d.logger != nil {
							d.logger.Error(
								"Failed to add refresh token JTI",
								slog.String("jti", issuedRefreshTokenJTI),
								slog.String("error", err.Error()),
							)
						}
					}
					for _, issuedAccessTokenJTI := range msg.IssuedAccessTokensJTIs {
						if _, err = db.Exec(
							InsertAccessTokenQuery,
							issuedAccessTokenJTI,
						); err != nil && d.logger != nil {
							d.logger.Error(
								"Failed to add access token JTI",
								slog.String("jti", issuedAccessTokenJTI),
								slog.String("error", err.Error()),
							)
						}
					}
					for _, revokedRefreshTokensJTI := range msg.RevokedRefreshTokensJTIs {
						if _, err = db.Exec(
							DeleteRefreshTokenQuery,
							revokedRefreshTokensJTI,
						); err != nil && d.logger != nil {
							d.logger.Error(
								"Failed to remove refresh token JTI",
								slog.String("jti", revokedRefreshTokensJTI),
								slog.String("error", err.Error()),
							)
						}
					}
					for _, revokedAccessTokensJTI := range msg.RevokedAccessTokensJTIs {
						if _, err = db.Exec(
							DeleteAccessTokenQuery,
							revokedAccessTokensJTI,
						); err != nil && d.logger != nil {
							d.logger.Error(
								"Failed to remove access token JTI",
								slog.String("jti", revokedAccessTokensJTI),
								slog.String("error", err.Error()),
							)
						}
					}
				}
			}
		},
	)

	// Wait for the goroutines to finish and return any errors
	err = eg.Wait()
	if err != nil && d.logger != nil {
		d.logger.Error(
			"Service encountered an error",
			slog.String("error", err.Error()),
		)
	}
	return err
}

// IsRefreshTokenValid checks if the given refresh token JTI exists in the database
//
// Parameters:
//
//   - jti: the refresh token JTI to validate
//
// Returns:
//
//   - bool: true if the refresh token JTI exists, false otherwise
//   - error: an error if the validation could not be performed
func (d *DefaultService) IsRefreshTokenValid(jti string) (bool, error) {
	// Check if the service is nil
	if d == nil {
		return false, sql.ErrConnDone
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return false, err
	}

	// Check if the refresh token JTI exists
	var exists bool
	err = db.QueryRow(
		CheckRefreshTokenQuery,
		jti,
	).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if d.logger != nil {
			d.logger.Error(
				"Failed to validate refresh token JTI",
				slog.String("jti", jti),
				slog.String("error", err.Error()),
			)
		}
		return false, err
	}
	return exists, nil
}

// IsAccessTokenValid checks if the given access token JTI exists in the database
//
// Parameters:
//
//   - jti: the access token JTI to validate
//
// Returns:
//
//   - bool: true if the access token JTI exists, false otherwise
//   - error: an error if the validation could not be performed
func (d *DefaultService) IsAccessTokenValid(jti string) (bool, error) {
	// Check if the service is nil
	if d == nil {
		return false, sql.ErrConnDone
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return false, err
	}

	// Check if the access token JTI exists
	var exists bool
	err = db.QueryRow(
		CheckAccessTokenQuery,
		jti,
	).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if d.logger != nil {
			d.logger.Error(
				"Failed to validate access token JTI",
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
func (d *DefaultService) ValidateClaims(
	claims map[string]interface{},
	token gojwttoken.Token,
) (bool, error) {
	// Check if the service is nil
	if d == nil {
		return false, sql.ErrConnDone
	}

	// Extract the JTI from the claims
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		if d.logger != nil {
			d.logger.Error(
				"JTI claim is missing or invalid",
				slog.String("token", token.String()),
			)
		}
		return false, nil
	}

	// Validate the JTI based on the token type
	switch token {
	case gojwttoken.AccessToken:
		return d.IsAccessTokenValid(jti)
	case gojwttoken.RefreshToken:
		return d.IsRefreshTokenValid(jti)
	default:
		if d.logger != nil {
			d.logger.Error(
				"Unknown token type",
				slog.String("token", token.String()),
			)
		}
		return false, nil
	}
}
