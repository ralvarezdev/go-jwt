package consumer

import (
	"context"
	"database/sql"
	"log/slog"

	gojwtrabbitmq "github.com/ralvarezdev/go-jwt/rabbitmq"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
	"golang.org/x/sync/errgroup"
)

type (
	// DefaultService is the default implementation of the Service interface
	DefaultService struct {
		gojwttokenclaims.TokenValidator
		logger   *slog.Logger
		consumer Consumer
	}
)

// NewDefaultService creates a new DefaultService
//
// Parameters:
//
//   - handler: the SQL connection handler
//   - consumer: the RabbitMQ consumer
//   - tokenValidator: the token validator
//   - logger: the logger (optional, can be nil)
//
// Returns:
//
//   - *DefaultService: the DefaultService instance
//   - error: an error if the data source or driver name is empty
func NewDefaultService(
	consumer Consumer,
	tokenValidator gojwttokenclaims.TokenValidator,
	logger *slog.Logger,
) (*DefaultService, error) {
	// Check if the consumer is nil
	if consumer == nil {
		return nil, gojwtrabbitmq.ErrNilConsumer
	}

	// Check if the token validator is nil
	if tokenValidator == nil {
		return nil, gojwttokenclaims.ErrNilClaimsValidator
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "jwt_rabbitmq_consumer_sql_service"),
		)
	}

	return &DefaultService{
		TokenValidator: tokenValidator,
		consumer:       consumer,
		logger:         logger,
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
						// Insert the refresh token ID
						if addErr := d.AddRefreshToken(
							ctx,
							issuedTokenPair.RefreshTokenID,
							issuedTokenPair.RefreshTokenExpiresAt,
						); addErr != nil {
							return addErr
						}

						// Also insert the access token ID
						if addErr := d.AddAccessToken(
							ctx,
							issuedTokenPair.AccessTokenID,
							issuedTokenPair.RefreshTokenID,
							issuedTokenPair.AccessTokenExpiresAt,
						); addErr != nil {
							return addErr
						}
					}

					// Remove the revoked refresh tokens ID
					for _, revokedTokenID := range msg.RevokedRefreshTokensID {
						if revokeErr := d.RevokeToken(
							ctx,
							gojwttoken.RefreshToken,
							revokedTokenID,
						); revokeErr != nil {
							return revokeErr
						}
					}

					// Remove the revoked refresh tokens ID
					for _, revokedTokenID := range msg.RevokedAccessTokensID {
						if revokeErr := d.RevokeToken(
							ctx,
							gojwttoken.AccessToken,
							revokedTokenID,
						); revokeErr != nil {
							return revokeErr
						}
					}

					// Remove the revoked access tokens ID
					for _, revokedTokenJTI := range msg.RevokedAccessTokensID {
						if revokeErr := d.RevokeToken(
							ctx,
							gojwttoken.AccessToken,
							revokedTokenJTI,
						); revokeErr != nil {
							return revokeErr
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
