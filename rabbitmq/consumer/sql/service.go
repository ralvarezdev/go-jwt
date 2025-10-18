package sql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	godatabases "github.com/ralvarezdev/go-databases"
	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
	gojwtrabbitmqconsumer "github.com/ralvarezdev/go-jwt/rabbitmq/consumer"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenvalidator "github.com/ralvarezdev/go-jwt/token/validator"
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
			slog.String("component", "jwt_rabbitmq_consumer_sql_service"),
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
//   - *sql.DB: the database connection
//   - error: an error if the connection could not be opened
func (d *DefaultService) Connect() (*sql.DB, error) {
	// Check if the service is nil
	if d == nil {
		return nil, godatabases.ErrNilService
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
		return nil, err
	}

	// Ensure the tables exist
	if _, err = db.Exec(CreateRefreshTokensTableQuery); err != nil {
		return nil, err
	}
	if _, err = db.Exec(CreateAccessTokensTableQuery); err != nil {
		return nil, err
	}
	return db, nil
}

// Disconnect closes the database connection
//
// Returns:
//
//   - error: an error if the connection could not be closed
func (d *DefaultService) Disconnect() error {
	// Check if the service is nil
	if d == nil {
		return godatabases.ErrNilService
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

// Start starts the service to listen for messages and update the SQL database
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

					// Check if the message is nil
					if msg == nil {
						if d.logger != nil {
							d.logger.Warn("Received nil message, skipping")
						}
						continue
					}

					// Process the message
					for _, issuedTokenPair := range msg.IssuedTokenPairs {
						// Insert the refresh token JTI
						if _, err = db.Exec(
							InsertRefreshTokenQuery,
							issuedTokenPair.RefreshTokenJTI,
						); err != nil && d.logger != nil {
							d.logger.Error(
								"Failed to add refresh token JTI",
								slog.String(
									"jti",
									issuedTokenPair.RefreshTokenJTI,
								),
								slog.String("error", err.Error()),
							)
						}

						// Also insert the access token JTI
						if _, err = db.Exec(
							InsertAccessTokenQuery,
							issuedTokenPair.AccessTokenJTI,
							issuedTokenPair.RefreshTokenJTI,
						); err != nil && d.logger != nil {
							d.logger.Error(
								"Failed to add access token JTI",
								slog.String(
									"jti",
									issuedTokenPair.AccessTokenJTI,
								),
								slog.String("error", err.Error()),
							)
						}
					}

					// Remove the revoked refresh token JTIs
					for _, revokedRefreshTokensJTI := range msg.RevokedRefreshTokensJTIs {
						// Remove the refresh token JTI
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

						// Also remove the associated access token JTIs
						if _, err = db.Exec(
							DeleteAccessTokenByRefreshTokenQuery,
							revokedRefreshTokensJTI,
						); err != nil && d.logger != nil {
							d.logger.Error(
								"Failed to remove access tokens by refresh token JTI",
								slog.String("jti", revokedRefreshTokensJTI),
								slog.String("error", err.Error()),
							)
						}
					}

					// Remove the revoked access token JTIs
					for _, revokedAccessTokensJTI := range msg.RevokedAccessTokensJTIs {
						// Remove the access token JTI
						if _, err = db.Exec(
							DeleteAccessTokenQuery,
							revokedAccessTokensJTI,
						); err != nil && d.logger != nil {
							d.logger.Error(
								"Failed to revoke access token JTI",
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
	claims jwt.MapClaims,
	token gojwttoken.Token,
) (bool, error) {
	// Check if the service is nil
	if d == nil {
		return false, sql.ErrConnDone
	}

	// Check if the claims are nil
	if claims == nil {
		if d.logger != nil {
			d.logger.Error(
				"Claims are nil",
				slog.String("token", token.String()),
			)
		}
		return false, gojwttokenvalidator.ErrNilClaims
	}

	// Check if the claims contain the JTI
	jti, ok := claims["jti"]
	if !ok {
		if d.logger != nil {
			d.logger.Error(
				"JTI claim is missing",
				slog.String("token", token.String()),
			)
		}
		return false, ErrMissingJTI
	}

	// Extract the JTI from the claims
	jtiStr, ok := jti.(string)
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
		return d.IsAccessTokenValid(jtiStr)
	case gojwttoken.RefreshToken:
		return d.IsRefreshTokenValid(jtiStr)
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
