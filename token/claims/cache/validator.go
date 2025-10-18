package cache

import (
	"log/slog"
	"time"

	gocache "github.com/ralvarezdev/go-cache"
	gocachetimed "github.com/ralvarezdev/go-cache/timed"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
)

type (
	// TokenValidatorService struct
	TokenValidatorService struct {
		logger *slog.Logger
		cache  *gocachetimed.Cache
	}
)

// NewTokenValidatorService creates a new token validator service
//
// Parameters:
//
//   - logger: The logger (optional, can be nil)
//
// Returns:
//
//   - *TokenValidatorService: The token validator service
func NewTokenValidatorService(logger *slog.Logger) *TokenValidatorService {
	if logger != nil {
		logger = logger.With(
			slog.String(
				"component",
				"token_validator_service",
			),
		)
	}
	return &TokenValidatorService{
		cache:  gocachetimed.NewCache(),
		logger: logger,
	}
}

// GetKey gets the key for the cache
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - string: The key for the cache
//   - error: An error if the token validator service is nil or if the token abbreviation fails
func (t *TokenValidatorService) GetKey(
	token gojwttoken.Token,
	id string,
) (string, error) {
	if t == nil {
		return "", gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the token string
	tokenPrefix, err := token.Abbreviation()
	if err != nil {
		return "", err
	}

	return gostringsadd.Prefixes(id, JwtIdentifierSeparator, tokenPrefix), nil
}

// Set sets a token in the cache
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//   - isValid: Whether the token is valid
//   - expiresAt: The expiration time of the token
//
// Returns:
//
//   - error: An error if the token validator service is nil or if setting the token in the cache fails
func (t *TokenValidatorService) Set(
	token gojwttoken.Token,
	id string,
	isValid bool,
	expiresAt time.Time,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return err
	}

	// Set the token in the cache
	err = t.cache.Set(key, gocachetimed.NewItem(isValid, expiresAt))
	if err != nil {
		gojwttokenclaims.SetTokenFailed(err, t.logger)
	}
	return err
}

// Revoke revokes a token in the cache
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - error: An error if the token validator service is nil or if revoking the token in the cache fails
func (t *TokenValidatorService) Revoke(
	token gojwttoken.Token,
	id string,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return err
	}

	// Revoke the token in the cache
	err = t.cache.UpdateValue(key, false)
	if err != nil {
		gojwttokenclaims.RevokeTokenFailed(err, t.logger)
	}
	return err
}

// IsValid checks if a token is valid in the cache
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - bool: Whether the token is valid
//   - error: An error if the token validator service is nil or if checking the token in the cache fails
func (t *TokenValidatorService) IsValid(
	token gojwttoken.Token,
	id string,
) (bool, error) {
	if t == nil {
		return false, gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return false, err
	}

	// Get the token from the cache
	isValid, found := t.cache.Get(key)
	if !found {
		return false, gocache.ErrItemNotFound
	}
	return isValid.(bool), nil
}
