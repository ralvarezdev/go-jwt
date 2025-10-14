package cache

import (
	"log/slog"

	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

// SetTokenToCache logs the set token to cache event
//
// Parameters:
//
//   - token: The token being set to cache
//   - id: The ID associated with the token
//   - logger: The logger to use for logging (optional, can be nil)
func SetTokenToCache(token gojwttoken.Token, id int64, logger *slog.Logger) {
	if logger != nil {
		logger.Debug(
			"Set token to cache",
			slog.String("token", token.String()),
			slog.Int64("id", id),
		)
	}
}

// SetTokenToCacheFailed logs the set token to cache failed event
//
// Parameters:
//
//   - err: The error that occurred while setting the token to cache
//   - logger: The logger to use for logging (optional, can be nil)
func SetTokenToCacheFailed(err error, logger *slog.Logger) {
	if logger != nil {
		logger.Error(
			"Set token to cache failed",
			slog.String("error", err.Error()),
		)
	}
}

// RevokeTokenFromCache logs the revoke token from cache event
//
// Parameters:
//
//   - token: The token being revoked from cache
//   - id: The ID associated with the token
//   - logger: The logger to use for logging (optional, can be nil)
func RevokeTokenFromCache(
	token gojwttoken.Token,
	id int64,
	logger *slog.Logger,
) {
	if logger != nil {
		logger.Debug(
			"Revoke token from cache",
			slog.String("token", token.String()),
			slog.Int64("id", id),
		)
	}
}

// RevokeTokenFromCacheFailed logs the revoke token from cache failed event
//
// Parameters:
//
//   - err: The error that occurred while revoking the token from cache
//   - logger: The logger to use for logging (optional, can be nil)
func RevokeTokenFromCacheFailed(err error, logger *slog.Logger) {
	if logger != nil {
		logger.Error(
			"Revoke token from cache failed",
			slog.String("error", err.Error()),
		)
	}
}

// GetTokenFromCache logs the get token from cache event
//
// Parameters:
//
//   - token: The token being retrieved from cache
//   - id: The ID associated with the token
//   - logger: The logger to use for logging (optional, can be nil)
func GetTokenFromCache(token gojwttoken.Token, id int64, logger *slog.Logger) {
	if logger != nil {
		logger.Debug(
			"Get token from cache",
			slog.String("token", token.String()),
			slog.Int64("id", id),
		)
	}
}

// GetTokenFromCacheFailed logs the get token from cache failed event
//
// Parameters:
//
//   - err: The error that occurred while retrieving the token from cache
//   - logger: The logger to use for logging (optional, can be nil)
func GetTokenFromCacheFailed(err error, logger *slog.Logger) {
	if logger != nil {
		logger.Error(
			"Get token from cache failed",
			slog.String("error", err.Error()),
		)
	}
}
