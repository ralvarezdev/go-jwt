package cache

import (
	"time"

	gocache "github.com/ralvarezdev/go-cache"
	gocachetimed "github.com/ralvarezdev/go-cache/timed"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
)

type (
	// TokenValidator interface
	TokenValidator interface {
		Set(
			token gojwttoken.Token,
			id string,
			isValid bool,
			expiresAt time.Time,
		) error
		Revoke(token gojwttoken.Token, id string) error
		IsValid(token gojwttoken.Token, id string) (bool, error)
	}

	// TokenValidatorService struct
	TokenValidatorService struct {
		logger *Logger
		cache  *gocachetimed.Cache
	}
)

// NewTokenValidatorService creates a new token validator service
func NewTokenValidatorService(logger *Logger) *TokenValidatorService {
	return &TokenValidatorService{
		cache:  gocachetimed.NewCache(),
		logger: logger,
	}
}

// GetKey gets the key for the cache
func (t *TokenValidatorService) GetKey(
	token gojwttoken.Token,
	id string,
) (string, error) {
	// Get the token string
	tokenPrefix, err := token.Abbreviation()
	if err != nil {
		return "", err
	}

	return gostringsadd.Prefixes(id, JwtIdentifierSeparator, tokenPrefix), nil
}

// Set sets a token in the cache
func (t *TokenValidatorService) Set(
	token gojwttoken.Token,
	id string,
	isValid bool,
	expiresAt time.Time,
) error {
	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return err
	}

	// Set the token in the cache
	err = t.cache.Set(key, gocachetimed.NewItem(isValid, expiresAt))
	if err != nil {
		// Log the error
		if t.logger != nil {
			t.logger.SetTokenToCacheFailed(err)
		}
	}
	return err
}

// Revoke revokes a token in the cache
func (t *TokenValidatorService) Revoke(
	token gojwttoken.Token,
	id string,
) error {
	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return err
	}

	// Revoke the token in the cache
	err = t.cache.UpdateValue(key, false)
	if err != nil {
		// Log the error
		if t.logger != nil {
			t.logger.RevokeTokenFromCacheFailed(err)
		}
	}
	return err
}

// IsValid checks if a token is valid in the cache
func (t *TokenValidatorService) IsValid(
	token gojwttoken.Token,
	id string,
) (bool, error) {
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
