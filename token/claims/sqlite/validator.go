package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	godatabases "github.com/ralvarezdev/go-databases"
	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
)

type (
	// TokenValidator is the default implementation of the Service interface
	TokenValidator struct {
		godatabasessql.Service
		logger *slog.Logger
	}
)

// NewTokenValidator creates a new TokenValidator
//
// Parameters:
//
//   - service: the SQL connection service
//   - logger: the logger (optional, can be nil)
//
// Returns:
//
//   - *TokenValidator: the TokenValidator instance
//   - error: an error if the data source or driver name is empty
func NewTokenValidator(
	service godatabasessql.Service,
	logger *slog.Logger,
) (*TokenValidator, error) {
	// Check if the service is nil
	if service == nil {
		return nil, godatabases.ErrNilService
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "sqlite_token_validator"),
		)
	}

	return &TokenValidator{
		Service: service,
		logger:  logger,
	}, nil
}

// Connect opens the database connection
//
// Parameters:
//
//   - ctx: the context
//
// Returns:
//
//   - *sql.DB: the database connection
//   - error: an error if the connection could not be opened
func (t *TokenValidator) Connect(ctx context.Context) (*sql.DB, error) {
	if t == nil {
		return nil, godatabases.ErrNilService
	}

	// Get the database connection
	db, err := t.Service.Connect()
	if err != nil {
		if t.logger != nil {
			t.logger.Error(
				"Failed to connect to database",
				slog.String("error", err.Error()),
			)
		}
		return nil, err
	}

	// Ensure the tables exist
	if _, err = db.ExecContext(ctx, CreateRefreshTokensTableQuery); err != nil {
		return nil, err
	}
	if _, err = db.ExecContext(ctx, CreateAccessTokensTableQuery); err != nil {
		return nil, err
	}
	return db, nil
}

// AddRefreshToken inserts a refresh token JTI into the database
//
// Parameters:
//
//   - ctx: the context for the query
//   - id: the refresh token JTI to insert
//   - expiresAt: the expiration time of the refresh token
//
// Returns:
//
//   - error: an error if the insertion could not be performed
func (t *TokenValidator) AddRefreshToken(
	ctx context.Context,
	id string,
	expiresAt time.Time,
) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Insert the refresh token JTI
	if _, err := t.ExecWithCtx(
		ctx,
		&InsertRefreshTokenQuery,
		id,
		expiresAt.Unix(),
	); err != nil && t.logger != nil {
		t.logger.Error(
			"Failed to insert refresh token JTI",
			slog.String("id", id),
			slog.String("error", err.Error()),
		)
	}
	return nil
}

// AddAccessToken inserts an access token JTI into the database
//
// Parameters:
//
//   - ctx: the context for the query
//   - id: the access token JTI to insert
//   - parentRefreshTokenID: the parent refresh token JTI
//   - expiresAt: the expiration time of the access token
//
// Returns:
//
//   - error: an error if the insertion could not be performed
func (t *TokenValidator) AddAccessToken(
	ctx context.Context,
	id, parentRefreshTokenID string,
	expiresAt time.Time,
) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Insert the access token JTI
	if _, err := t.ExecWithCtx(
		ctx,
		&InsertAccessTokenQuery,
		id,
		parentRefreshTokenID,
		expiresAt.Unix(),
	); err != nil && t.logger != nil {
		t.logger.Error(
			"Failed to insert access token JTI",
			slog.String("id", id),
			slog.String("error", err.Error()),
		)
	}
	return nil
}

// RevokeAccessTokenByRefreshToken revokes access tokens associated with the given refresh token JTIs
//
// Parameters:
//
//   - ctx: the context for the query
//   - id: the refresh token JTI whose associated access tokens are to be revoked
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (t *TokenValidator) RevokeAccessTokenByRefreshToken(ctx context.Context, id string) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Revoke the access tokens associated with the refresh token JTI
	if _, err := t.ExecWithCtx(
		ctx,
		&DeleteAccessTokenByRefreshTokenQuery,
		id,
	); err != nil && t.logger != nil {
		t.logger.Error(
			"Failed to revoke access tokens by refresh token JTI",
			slog.String("id", id),
			slog.String("error", err.Error()),
		)
	}
	return nil
}

// RevokeRefreshToken revokes a refresh token JTI from the database
//
// Parameters:
//
//   - ctx: the context for the query
//   - id: the refresh token JTI to revoke
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (t *TokenValidator) RevokeRefreshToken(ctx context.Context, id string) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Revoke the refresh token JTI
	if _, err := t.ExecWithCtx(
		ctx,
		&DeleteRefreshTokenQuery,
		id,
	); err != nil && t.logger != nil {
		t.logger.Error(
			"Failed to revoke refresh token JTI",
			slog.String("id", id),
			slog.String("error", err.Error()),
		)
	}

	// Revoke associated access tokens first
	return t.RevokeAccessTokenByRefreshToken(ctx, id)
}

// RevokeAccessToken revokes an access token JTI from the database
//
// Parameters:
//
//   - ctx: the context for the query
//   - id: the access token JTI to revoke
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (t *TokenValidator) RevokeAccessToken(ctx context.Context, id string) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Revoke the access token JTI
	if _, err := t.ExecWithCtx(
		ctx,
		&DeleteAccessTokenQuery,
		id,
	); err != nil && t.logger != nil {
		t.logger.Error(
			"Failed to revoke access token JTI",
			slog.String("id", id),
			slog.String("error", err.Error()),
		)
	}
	return nil
}

// RevokeToken revokes a token JTI from the database based on the token type
//
// Parameters:
//
//   - ctx: the context for the query
//   - token: the token type (access or refresh)
//   - id: the token JTI to revoke
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (t *TokenValidator) RevokeToken(
	ctx context.Context,
	token gojwttoken.Token,
	id string,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Revoke the JTI based on the token type
	switch token {
	case gojwttoken.AccessToken:
		return t.RevokeAccessToken(ctx, id)
	case gojwttoken.RefreshToken:
		return t.RevokeRefreshToken(ctx, id)
	default:
		if t.logger != nil {
			t.logger.Error(
				"Unknown token type",
				slog.String("token", token.String()),
			)
		}
		return nil
	}
}

// IsRefreshTokenValid checks if the given refresh token JTI exists in the database
//
// Parameters:
//
//   - id: the refresh token JTI to validate
//
// Returns:
//
//   - bool: true if the refresh token JTI exists, false otherwise
//   - error: an error if the validation could not be performed
func (t *TokenValidator) IsRefreshTokenValid(ctx context.Context, id string) (bool, error) {
	// Check if the service is nil
	if t == nil {
		return false, gojwttokenclaims.ErrNilTokenValidator
	}

	// Check if the refresh token JTI exists
	return t.IsTokenValid(ctx, gojwttoken.RefreshToken, id)
}

// IsAccessTokenValid checks if the given access token JTI exists in the database
//
// Parameters:
//
//   - ctx: the context for the query
//   - id: the access token JTI to validate
//
// Returns:
//
//   - bool: true if the access token JTI exists, false otherwise
//   - error: an error if the validation could not be performed
func (t *TokenValidator) IsAccessTokenValid(ctx context.Context, id string) (bool, error) {
	// Check if the service is nil
	if t == nil {
		return false, gojwttokenclaims.ErrNilTokenValidator
	}

	// Check if the access token JTI exists
	return t.IsTokenValid(ctx, gojwttoken.AccessToken, id)
}

// IsTokenValid validates the token
//
// Parameters:
//
//   - ctx: the context for the query
//   - token: the token type
//   - id: the ID associated with the token
//
// Returns:
//
//   - bool: true if the claims are valid, false otherwise
//   - error: an error if the validation could not be performed
func (t *TokenValidator) IsTokenValid(ctx context.Context, token gojwttoken.Token, id string) (
	bool,
	error,
) {
	// Check if the service is nil
	if t == nil {
		return false, gojwttokenclaims.ErrNilTokenValidator
	}

	// Determine the query based on the token type
	var query string
	switch token {
	case gojwttoken.AccessToken:
		query = CheckAccessTokenQuery
	case gojwttoken.RefreshToken:
		query = CheckRefreshTokenQuery
	default:
		if t.logger != nil {
			t.logger.Error(
				"Unknown token type",
				slog.String("token", token.String()),
			)
		}
		return false, nil
	}

	// Check if the refresh token JTI exists
	var exists bool
	rows, err := t.QueryRowWithCtx(
		ctx,
		&query,
		id,
	)
	if err != nil {
		if t.logger != nil && token == gojwttoken.AccessToken {
			t.logger.Error(
				"Failed to validate access token JTI",
				slog.String("id", id),
				slog.String("error", err.Error()),
			)
		}
		if t.logger != nil && token == gojwttoken.RefreshToken {
			t.logger.Error(
				"Failed to validate refresh token JTI",
				slog.String("id", id),
				slog.String("error", err.Error()),
			)
		}
		return false, err
	}
	if rowsErr := rows.Scan(&exists); rowsErr != nil && !errors.Is(rowsErr, sql.ErrNoRows) {
		if t.logger != nil && token == gojwttoken.AccessToken {
			t.logger.Error(
				"Failed to validate access token JTI",
				slog.String("id", id),
				slog.String("error", rowsErr.Error()),
			)
		} else if t.logger != nil && token == gojwttoken.RefreshToken {
			t.logger.Error(
				"Failed to validate refresh token JTI",
				slog.String("id", id),
				slog.String("error", rowsErr.Error()),
			)
		}
		return false, rowsErr
	}
	return exists, nil
}
