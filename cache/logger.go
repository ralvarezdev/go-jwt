package cache

import (
	"fmt"

	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gologgermode "github.com/ralvarezdev/go-logger/mode"
	gologgermodenamed "github.com/ralvarezdev/go-logger/mode/named"
)

// Logger is the cache token validator logger
type Logger struct {
	logger gologgermodenamed.Logger
}

// NewLogger creates a new cache token validator logger
func NewLogger(header string, modeLogger gologgermode.Logger) (*Logger, error) {
	// Initialize the mode named logger
	namedLogger, err := gologgermodenamed.NewDefaultLogger(header, modeLogger)
	if err != nil {
		return nil, err
	}

	return &Logger{logger: namedLogger}, nil
}

// SetTokenToCache logs the set token to cache event
func (l *Logger) SetTokenToCache(token gojwttoken.Token, id int64) {
	l.logger.Debug(
		"set token to cache",
		fmt.Sprintf("token: %s, id: %d", token, id),
	)
}

// SetTokenToCacheFailed logs the set token to cache failed event
func (l *Logger) SetTokenToCacheFailed(err error) {
	l.logger.Error(
		"set token to cache failed",
		err,
	)
}

// RevokeTokenFromCache logs the revoke token from cache event
func (l *Logger) RevokeTokenFromCache(token gojwttoken.Token, id int64) {
	l.logger.Debug(
		"revoke token from cache",
		fmt.Sprintf("token: %s, id: %d", token, id),
	)
}

// RevokeTokenFromCacheFailed logs the revoke token from cache failed event
func (l *Logger) RevokeTokenFromCacheFailed(err error) {
	l.logger.Error(
		"revoke token from cache failed",
		err,
	)
}

// GetTokenFromCache logs the get token from cache event
func (l *Logger) GetTokenFromCache(token gojwttoken.Token, id int64) {
	l.logger.Debug(
		"get token from cache",
		fmt.Sprintf("token: %s, id: %d", token, id),
	)
}

// GetTokenFromCacheFailed logs the get token from cache failed event
func (l *Logger) GetTokenFromCacheFailed(err error) {
	l.logger.Error(
		"get token from cache failed",
		err,
	)
}
