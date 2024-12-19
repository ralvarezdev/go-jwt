package validator

import (
	gologger "github.com/ralvarezdev/go-logger"
	gologgerstatus "github.com/ralvarezdev/go-logger/status"
)

// Logger is the JWT validator logger
type Logger struct {
	logger gologger.Logger
}

// NewLogger creates a new JWT validator logger
func NewLogger(logger gologger.Logger) (*Logger, error) {
	// Check if the logger is nil
	if logger == nil {
		return nil, gologger.NilLoggerError
	}

	return &Logger{logger: logger}, nil
}

// ValidatedToken logs a message when the server validates a token
func (l *Logger) ValidatedToken() {
	l.logger.LogMessage(
		gologger.NewLogMessage(
			"Validated token",
			gologgerstatus.StatusInfo,
			nil,
		),
	)
}

// MissingTokenClaimsUserId logs the missing token claims user ID
func (l *Logger) MissingTokenClaimsUserId() {
	l.logger.LogMessage(
		gologger.NewLogMessage(
			"Missing  user ID in token claims",
			gologgerstatus.StatusFailed,
			nil,
		),
	)
}
