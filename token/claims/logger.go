package claims

import (
	"log/slog"

	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

// SetToken logs the set token
//
// Parameters:
//
//   - token: The token being set
//   - id: The ID associated with the token
//   - logger: The logger to use for logging (optional, can be nil)
func SetToken(token gojwttoken.Token, id int64, logger *slog.Logger) {
	if logger != nil {
		logger.Debug(
			"Set token",
			slog.String("token", token.String()),
			slog.Int64("id", id),
		)
	}
}

// SetTokenFailed logs the set token failure
//
// Parameters:
//
//   - err: The error that occurred while setting the token
//   - logger: The logger to use for logging (optional, can be nil)
func SetTokenFailed(err error, logger *slog.Logger) {
	if logger != nil {
		logger.Error(
			"Set token failed",
			slog.String("error", err.Error()),
		)
	}
}

// RevokeToken logs the revoke token
//
// Parameters:
//
//   - token: The token being revoked
//   - id: The ID associated with the token
//   - logger: The logger to use for logging (optional, can be nil)
func RevokeToken(
	token gojwttoken.Token,
	id int64,
	logger *slog.Logger,
) {
	if logger != nil {
		logger.Debug(
			"Revoke token",
			slog.String("token", token.String()),
			slog.Int64("id", id),
		)
	}
}

// RevokeTokenFailed logs the revoke token failure
//
// Parameters:
//
//   - err: The error that occurred while revoking the token
//   - logger: The logger to use for logging (optional, can be nil)
func RevokeTokenFailed(err error, logger *slog.Logger) {
	if logger != nil {
		logger.Error(
			"Revoke token failed",
			slog.String("error", err.Error()),
		)
	}
}

// GetToken logs the get token
//
// Parameters:
//
//   - token: The token being retrieved
//   - id: The ID associated with the token
//   - logger: The logger to use for logging (optional, can be nil)
func GetToken(token gojwttoken.Token, id int64, logger *slog.Logger) {
	if logger != nil {
		logger.Debug(
			"Get token",
			slog.String("token", token.String()),
			slog.Int64("id", id),
		)
	}
}

// GetTokenFailed logs the get token failure
//
// Parameters:
//
//   - err: The error that occurred while retrieving the token
//   - logger: The logger to use for logging (optional, can be nil)
func GetTokenFailed(err error, logger *slog.Logger) {
	if logger != nil {
		logger.Error(
			"Get token failed",
			slog.String("error", err.Error()),
		)
	}
}
