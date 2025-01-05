package issuer

import (
	gologgermode "github.com/ralvarezdev/go-logger/mode"
	gologgermodenamed "github.com/ralvarezdev/go-logger/mode/named"
)

// Logger is the JWT issuer logger
type Logger struct {
	logger gologgermodenamed.Logger
}

// NewLogger creates a new JWT issuer logger
func NewLogger(header string, modeLogger gologgermode.Logger) (*Logger, error) {
	// Initialize the mode named logger
	namedLogger, err := gologgermodenamed.NewDefaultLogger(header, modeLogger)
	if err != nil {
		return nil, err
	}

	return &Logger{logger: namedLogger}, nil
}

// IssuedToken logs a message when the server issues a token
func (l *Logger) IssuedToken() {
	l.logger.Debug("issued token")
}
