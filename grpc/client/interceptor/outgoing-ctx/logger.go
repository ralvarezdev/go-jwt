package outgoing_ctx

import (
	gologger "github.com/ralvarezdev/go-logger"
	gologgerstatus "github.com/ralvarezdev/go-logger/status"
)

// Logger is the logger for the outgoing context debugger
type Logger struct {
	logger gologger.Logger
}

// NewLogger is the logger for the outgoing context debugger
func NewLogger(logger gologger.Logger) (*Logger, error) {
	// Check if the logger is nil
	if logger == nil {
		return nil, gologger.NilLoggerError
	}

	return &Logger{logger: logger}, nil
}

// LogKeyValue logs the key value
func (l *Logger) LogKeyValue(key string, value string) {
	formattedKey := "Outgoing context key '" + key + "' value"
	l.logger.LogMessage(gologger.NewLogMessage(formattedKey, gologgerstatus.StatusDebug, nil, value))
}
