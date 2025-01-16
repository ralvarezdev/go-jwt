package cache

import (
	gocachetimed "github.com/ralvarezdev/go-cache/timed"
	"time"
)

type (
	// TokenValidator interface
	TokenValidator interface {
		Set(id string, value interface{}, period time.Duration) error
		Has(id string) (bool, error)
		Get(id string) (interface{}, bool)
		Delete(id string) error
	}

	// TokenValidatorService struct
	TokenValidatorService struct {
		gocachetimed.Cache
	}
)

// NewTokenValidatorService creates a new token validator service
func NewTokenValidatorService() *TokenValidatorService {
	return &TokenValidatorService{
		Cache: gocachetimed.Cache{},
	}
}

// Set sets a token in the cache
func (t *TokenValidatorService) Set(
	id string,
	value interface{},
	period time.Duration,
) error {
	return t.Cache.Set(id, gocachetimed.NewItem(value, time.Now().Add(period)))
}

// Has checks if a token exists in the cache
func (t *TokenValidatorService) Has(id string) (bool, error) {
	return t.Cache.Has(id), nil
}

// Get gets a token from the cache
func (t *TokenValidatorService) Get(id string) (interface{}, bool) {
	return t.Cache.Get(id)
}

// Delete deletes a token from the cache
func (t *TokenValidatorService) Delete(id string) error {
	t.Cache.Delete(id)
	return nil
}
