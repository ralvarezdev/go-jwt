package sqlite

import (
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"time"

	godatabases "github.com/ralvarezdev/go-databases"
	godatabasessql "github.com/ralvarezdev/go-databases/sql"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
	gojwttokenvalidator "github.com/ralvarezdev/go-jwt/token/validator"
)

type (
	// TokenValidator is the default implementation of the Service interface
	TokenValidator struct {
		godatabasessql.Handler
		logger *slog.Logger
		mutex  sync.Mutex
	}
)

// NewTokenValidator creates a new TokenValidator
//
// Parameters:
//
//   - handler: the SQL connection handler
//   - logger: the logger (optional, can be nil)
//
// Returns:
//
//   - *TokenValidator: the TokenValidator instance
//   - error: an error if the data source or driver name is empty
func NewTokenValidator(
	handler godatabasessql.Handler,
	logger *slog.Logger,
) (*TokenValidator, error) {
	// Check if the handler is nil
	if handler == nil {
		return nil, godatabases.ErrNilHandler
	}

	if logger != nil {
		logger = logger.With(
			slog.String("component", "sqlite_token_validator"),
		)
	}

	return &TokenValidator{
		Handler: handler,
		logger:  logger,
	}, nil
}

// Connect opens the database connection
//
// Returns:
//
//   - *sql.DB: the database connection
//   - error: an error if the connection could not be opened
func (t *TokenValidator) Connect() (*sql.DB, error) {
	if t == nil {
		return nil, godatabases.ErrNilService
	}

	// Lock the mutex to ensure thread safety
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Connect to the database
	db, err := t.Handler.Connect()
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
func (t *TokenValidator) Disconnect() error {
	if t == nil {
		return gojwttokenvalidator.ErrNilValidator
	}

	// Lock the mutex to ensure thread safety
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Disconnect from the database
	if err := t.Handler.Disconnect(); err != nil {
		if t.logger != nil {
			t.logger.Error(
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
func (t *TokenValidator) DB() (*sql.DB, error) {
	// Lock the mutex to ensure thread safety
	t.mutex.Lock()

	// Get the database connection
	db, err := t.DB()
	if err != nil {
		t.mutex.Unlock()
		if t.logger != nil {
			t.logger.Error(
				"Failed to get database connection",
				slog.String("error", err.Error()),
			)
		}
		return nil, err
	}
	t.mutex.Unlock()

	return db, nil
}

// AddRefreshToken inserts a refresh token JTI into the database
//
// Parameters:
//
//   - id: the refresh token JTI to insert
//   - expiresAt: the expiration time of the refresh token
//
// Returns:
//
//   - error: an error if the insertion could not be performed
func (t *TokenValidator) AddRefreshToken(
	id string,
	expiresAt time.Time,
) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the database connection
	db, err := t.DB()
	if err != nil {
		return err
	}

	// Insert the refresh token JTI
	if _, err = db.Exec(
		InsertRefreshTokenQuery,
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
//   - id: the access token JTI to insert
//   - parentRefreshTokenID: the parent refresh token JTI
//   - expiresAt: the expiration time of the access token
//
// Returns:
//
//   - error: an error if the insertion could not be performed
func (t *TokenValidator) AddAccessToken(
	id, parentRefreshTokenID string,
	expiresAt time.Time,
) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the database connection
	db, err := t.DB()
	if err != nil {
		return err
	}

	// Insert the access token JTI
	if _, err = db.Exec(
		InsertAccessTokenQuery,
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
//   - id: the refresh token JTI whose associated access tokens are to be revoked
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (t *TokenValidator) RevokeAccessTokenByRefreshToken(id string) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the database connection
	db, err := t.DB()
	if err != nil {
		return err
	}

	// Revoke the access tokens associated with the refresh token JTI
	if _, err = db.Exec(
		DeleteAccessTokenByRefreshTokenQuery,
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
//   - id: the refresh token JTI to revoke
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (t *TokenValidator) RevokeRefreshToken(id string) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the database connection
	db, err := t.DB()
	if err != nil {
		return err
	}

	// Revoke the refresh token JTI
	if _, err = db.Exec(
		DeleteRefreshTokenQuery,
		id,
	); err != nil && t.logger != nil {
		t.logger.Error(
			"Failed to revoke refresh token JTI",
			slog.String("id", id),
			slog.String("error", err.Error()),
		)
	}

	// Revoke associated access tokens first
	if err = t.RevokeAccessTokenByRefreshToken(id); err != nil {
		return err
	}
	return nil
}

// RevokeAccessToken revokes an access token JTI from the database
//
// Parameters:
//
//   - id: the access token JTI to revoke
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (t *TokenValidator) RevokeAccessToken(id string) error {
	// Check if the service is nil
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the database connection
	db, err := t.DB()
	if err != nil {
		return err
	}

	// Revoke the access token JTI
	if _, err = db.Exec(
		DeleteAccessTokenQuery,
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
//   - token: the token type (access or refresh)
//   - id: the token JTI to revoke
//
// Returns:
//
//   - error: an error if the revocation could not be performed
func (t *TokenValidator) RevokeToken(
	token gojwttoken.Token,
	id string,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Revoke the JTI based on the token type
	switch token {
	case gojwttoken.AccessToken:
		return t.RevokeAccessToken(id)
	case gojwttoken.RefreshToken:
		return t.RevokeRefreshToken(id)
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
func (t *TokenValidator) IsRefreshTokenValid(id string) (bool, error) {
	// Check if the service is nil
	if t == nil {
		return false, gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the database connection
	db, err := t.DB()
	if err != nil {
		return false, err
	}

	// Check if the refresh token JTI exists
	var exists bool
	err = db.QueryRow(
		CheckRefreshTokenQuery,
		id,
	).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if t.logger != nil {
			t.logger.Error(
				"Failed to validate refresh token JTI",
				slog.String("id", id),
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
//   - id: the access token JTI to validate
//
// Returns:
//
//   - bool: true if the access token JTI exists, false otherwise
//   - error: an error if the validation could not be performed
func (t *TokenValidator) IsAccessTokenValid(id string) (bool, error) {
	// Check if the service is nil
	if t == nil {
		return false, gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the database connection
	db, err := t.DB()
	if err != nil {
		return false, err
	}

	// Check if the access token JTI exists
	var exists bool
	err = db.QueryRow(
		CheckAccessTokenQuery,
		id,
	).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if t.logger != nil {
			t.logger.Error(
				"Failed to validate access token JTI",
				slog.String("id", id),
				slog.String("error", err.Error()),
			)
		}
		return false, err
	}
	return exists, nil
}

// IsTokenValid validates the token
//
// Parameters:
//
//   - token: the token type
//   - id: the ID associated with the token
//
// Returns:
//
//   - bool: true if the claims are valid, false otherwise
//   - error: an error if the validation could not be performed
func (t *TokenValidator) IsTokenValid(token gojwttoken.Token, id string) (
	bool,
	error,
) {
	if t == nil {
		return false, gojwttokenclaims.ErrNilTokenValidator
	}

	// Validate the JTI based on the token type
	switch token {
	case gojwttoken.AccessToken:
		return t.IsAccessTokenValid(id)
	case gojwttoken.RefreshToken:
		return t.IsRefreshTokenValid(id)
	default:
		if t.logger != nil {
			t.logger.Error(
				"Unknown token type",
				slog.String("token", token.String()),
			)
		}
		return false, nil
	}
}
