package cache

import (
	gocachetimed "github.com/ralvarezdev/go-cache/timed"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
	gostringsseparator "github.com/ralvarezdev/go-strings/separator"
	"time"
)

type (
	// TokenValidator interface
	TokenValidator interface {
		Set(
			token gojwttoken.Token,
			id string,
			value interface{},
			expiresAt time.Time,
		) error
		Has(token gojwttoken.Token, id string) (bool, error)
		Get(token gojwttoken.Token, id string) (interface{}, bool)
		Delete(token gojwttoken.Token, id string) error
	}

	// TokenValidatorService struct
	TokenValidatorService struct {
		logger *Logger
		gocachetimed.Cache
	}
)

// NewTokenValidatorService creates a new token validator service
func NewTokenValidatorService(logger *Logger) *TokenValidatorService {
	return &TokenValidatorService{
		Cache:  gocachetimed.Cache{},
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

	return gostringsadd.Prefixes(id, gostringsseparator.Dots, tokenPrefix), nil
}

// Set sets a token in the cache
func (t *TokenValidatorService) Set(
	token gojwttoken.Token,
	id string,
	value interface{},
	expiresAt time.Time,
) error {
	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return err
	}

	// Set the token in the cache
	err = t.Cache.Set(key, gocachetimed.NewItem(value, expiresAt))
	if err != nil {
		// Log the error
		if t.logger != nil {
			t.logger.SetTokenToCacheFailed(err)
		}
	}
	return err
}

// Has checks if a token exists in the cache
func (t *TokenValidatorService) Has(
	token gojwttoken.Token,
	id string,
) (bool, error) {
	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return false, err
	}

	// Check if the token exists in the cache
	return t.Cache.Has(key), nil
}

// Get gets a token from the cache
func (t *TokenValidatorService) Get(
	token gojwttoken.Token,
	id string,
) (interface{}, bool) {
	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return nil, false
	}

	// Get the token from the cache
	return t.Cache.Get(key)
}

// Delete deletes a token from the cache
func (t *TokenValidatorService) Delete(
	token gojwttoken.Token,
	id string,
) error {
	// Get the key
	key, err := t.GetKey(token, id)
	if err != nil {
		return err
	}

	// Delete the token from the cache
	t.Cache.Delete(key)
	return nil
}
