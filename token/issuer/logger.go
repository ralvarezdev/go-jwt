package issuer

import (
	gologger "github.com/ralvarezdev/go-logger"
	gologgerstatus "github.com/ralvarezdev/go-logger/status"
)

// Logger is the JWT issuer logger
type Logger struct {
	logger gologger.Logger
}

// NewLogger creates a new JWT issuer logger
func NewLogger(logger gologger.Logger) (*Logger, error) {
	// Check if the logger is nil
	if logger == nil {
		return nil, gologger.NilLoggerError
	}

	return &Logger{logger: logger}, nil
}

// IssuedToken logs a message when the server issues a token
func (l *Logger) IssuedToken() {
	l.logger.LogMessage(gologger.NewLogMessage("issued token", gologgerstatus.StatusInfo, nil))
}
