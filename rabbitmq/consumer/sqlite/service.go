package sqlite

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
						if err = d.InsertRefreshTokens(issuedTokenPair.RefreshTokenJTI); err != nil {
							return err
						}

						// Also insert the access token JTI
						if err = d.InsertAccessTokens(issuedTokenPair); err != nil {
							return err
						}
					}

					// Remove the revoked refresh token JTIs
					if err = d.RevokeRefreshTokens(msg.RevokedRefreshTokensJTIs...); err != nil {
						return err
					}

					// Remove the revoked refresh token JTIs
					if err = d.RevokeAccessTokensByRefreshTokens(msg.RevokedRefreshTokensJTIs...); err != nil {
						return err
					}

					// Remove the revoked access token JTIs
					if err = d.RevokeAccessTokens(msg.RevokedAccessTokensJTIs...); err != nil {
						return err
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

// InsertRefreshTokens inserts multiple refresh token JTIs into the database
//
// Parameters:
//
//   - jtis: the slice of refresh token JTIs to insert
//
// Returns:
//
//   - error: an error if the insertion could not be performed
func (d *DefaultService) InsertRefreshTokens(jtis ...string) error {
	// Check if the service is nil
	if d == nil {
		return sql.ErrConnDone
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return err
	}

	// Insert each refresh token JTI
	for _, jti := range jtis {
		if _, err = db.Exec(
			InsertRefreshTokenQuery,
			jti,
		); err != nil && d.logger != nil {
			d.logger.Error(
				"Failed to insert refresh token JTI",
				slog.String("jti", jti),
				slog.String("error", err.Error()),
			)
		}
	}
	return nil
}

// InsertAccessTokens inserts multiple access token JTIs into the database
//
// Parameters:
//
//   - tokenPairs: the slice of access token JTI pairs to insert
//
// Returns:
//
//   - error: an error if the insertion could not be performed
func (d *DefaultService) InsertAccessTokens(tokenPairs ...gojwtrabbitmq.TokenPair) error {
	// Check if the service is nil
	if d == nil {
		return sql.ErrConnDone
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return err
	}

	// Insert each access token JTI
	for _, tokenPair := range tokenPairs {
		if _, err = db.Exec(
			InsertAccessTokenQuery,
			tokenPair.AccessTokenJTI,
			tokenPair.RefreshTokenJTI,
		); err != nil && d.logger != nil {
			d.logger.Error(
				"Failed to insert access token JTI",
				slog.String("jti", tokenPair.AccessTokenJTI),
				slog.String("error", err.Error()),
			)
		}
	}
	return nil
}

// RevokeRefreshTokens revokes multiple refresh token JTIs from the database
//
// Parameters:
//
//   - jtis: the slice of refresh token JTIs to revoke
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (d *DefaultService) RevokeRefreshTokens(jtis ...string) error {
	// Check if the service is nil
	if d == nil {
		return sql.ErrConnDone
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return err
	}

	// Revoke each refresh token JTI
	for _, jti := range jtis {
		if _, err = db.Exec(
			DeleteRefreshTokenQuery,
			jti,
		); err != nil && d.logger != nil {
			d.logger.Error(
				"Failed to revoke refresh token JTI",
				slog.String("jti", jti),
				slog.String("error", err.Error()),
			)
		}
	}
	return nil
}

// RevokeAccessTokens revokes multiple access token JTIs from the database
//
// Parameters:
//
//   - jtis: the slice of access token JTIs to revoke
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (d *DefaultService) RevokeAccessTokens(jtis ...string) error {
	// Check if the service is nil
	if d == nil {
		return sql.ErrConnDone
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return err
	}

	// Revoke each access token JTI
	for _, jti := range jtis {
		if _, err = db.Exec(
			DeleteAccessTokenQuery,
			jti,
		); err != nil && d.logger != nil {
			d.logger.Error(
				"Failed to revoke access token JTI",
				slog.String("jti", jti),
				slog.String("error", err.Error()),
			)
		}
	}
	return nil
}

// RevokeAccessTokensByRefreshTokens revokes access token JTIs associated with the given refresh tokens JTIs
//
// Parameters:
//
//   - jtis: the slice of refresh token JTIs whose associated access tokens are to be revoked
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (d *DefaultService) RevokeAccessTokensByRefreshTokens(jtis ...string) error {
	// Check if the service is nil
	if d == nil {
		return sql.ErrConnDone
	}

	// Get the database connection
	db, err := d.DB()
	if err != nil {
		return err
	}

	// Revoke access tokens for each refresh token JTI
	for _, jti := range jtis {
		if _, err = db.Exec(
			DeleteAccessTokenByRefreshTokenQuery,
			jti,
		); err != nil && d.logger != nil {
			d.logger.Error(
				"Failed to revoke access tokens by refresh token JTI",
				slog.String("jti", jti),
				slog.String("error", err.Error()),
			)
		}
	}
	return nil
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

/*
Set(
			token gojwttoken.Token,
			id string,
			isValid bool,
			expiresAt time.Time,
		) error
		Revoke(token gojwttoken.Token, id string) error
		IsValid(token gojwttoken.Token, id string) (bool, error)
*/
