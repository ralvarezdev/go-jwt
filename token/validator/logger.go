package validator

import (
	gologgermode "github.com/ralvarezdev/go-logger/mode"
	gologgermodenamed "github.com/ralvarezdev/go-logger/mode/named"
)

// Logger is the JWT validator logger
type Logger struct {
	logger gologgermodenamed.Logger
}

// NewLogger creates a new JWT validator logger
func NewLogger(header string, modeLogger gologgermode.Logger) (*Logger, error) {
	// Initialize the mode named logger
	namedLogger, err := gologgermodenamed.NewDefaultLogger(header, modeLogger)
	if err != nil {
		return nil, err
	}

	return &Logger{logger: namedLogger}, nil
}

// ValidatedToken logs a message when the server validates a token
func (l *Logger) ValidatedToken() {
	l.logger.Debug("validated token")
}

// MissingTokenClaimsUserId logs the missing token claims user ID
func (l *Logger) MissingTokenClaimsUserId() {
	l.logger.Warning("missing  user id in token claims")
}
