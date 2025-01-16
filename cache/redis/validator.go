package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	godatabasesredis "github.com/ralvarezdev/go-databases/redis"
	gojwtcache "github.com/ralvarezdev/go-jwt/cache"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
	gostringsseparator "github.com/ralvarezdev/go-strings/separator"
	"time"
)

type (
	// TokenValidatorService struct
	TokenValidatorService struct {
		redisClient *redis.Client
		logger      *gojwtcache.Logger
	}
)

// NewTokenValidatorService creates a new token validator service
func NewTokenValidatorService(
	redisClient *redis.Client,
	logger *gojwtcache.Logger,
) (
	*TokenValidatorService,
	error,
) {
	// Check if the Redis client is nil
	if redisClient == nil {
		return nil, godatabasesredis.ErrNilClient
	}

	return &TokenValidatorService{redisClient: redisClient, logger: logger}, nil
}

// GetKey gets the JWT Identifier key
func (d *TokenValidatorService) GetKey(
	token gojwttoken.Token,
	id string,
) (string, error) {
	// Get the token string
	tokenPrefix, err := token.Abbreviation()
	if err != nil {
		return "", err
	}

	return gostringsadd.Prefixes(
		id,
		gostringsseparator.Dots,
		JwtIdentifierPrefix,
		tokenPrefix,
	), nil
}

// Set sets the token with the value and period
func (d *TokenValidatorService) Set(
	token gojwttoken.Token,
	id string,
	value interface{},
	expiresAt time.Time,
) error {
	// Get the key
	key, err := d.GetKey(token, id)
	if err != nil {
		return err
	}

	// Set the initial value
	if err = d.redisClient.Set(
		context.Background(),
		key,
		value,
		0,
	).Err(); err != nil {
		// Log the error
		if d.logger != nil {
			d.logger.SetTokenToCacheFailed(err)
		}
		return err
	}

	// Set expiration time for the key as a UNIX timestamp
	err = d.redisClient.ExpireAt(context.Background(), key, expiresAt).Err()
	if err != nil {
		// Log the error
		if d.logger != nil {
			d.logger.SetTokenToCacheFailed(err)
		}
	}
	return err
}

// Has checks if the token is valid
func (d *TokenValidatorService) Has(
	token gojwttoken.Token,
	id string,
) (bool, error) {
	// Get the key
	key, err := d.GetKey(token, id)
	if err != nil {
		return false, err
	}

	// Check the JWT Identifier
	_, err = d.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		// Log the error
		if d.logger != nil {
			d.logger.HasTokenInCacheFailed(err)
		}
		return false, err
	}
	return true, nil
}

// Get gets the token
func (d *TokenValidatorService) Get(
	token gojwttoken.Token,
	id string,
) (interface{}, error) {
	// Get the key
	key, err := d.GetKey(token, id)
	if err != nil {
		return nil, err
	}

	// Get the value
	value, err := d.redisClient.Get(
		context.Background(),
		key,
	).Result()
	if err != nil {
		// Log the error
		if d.logger != nil {
			d.logger.GetTokenFromCacheFailed(err)
		}
		return nil, err
	}
	return value, err
}

// Delete deletes the token
func (d *TokenValidatorService) Delete(
	token gojwttoken.Token,
	id string,
) error {
	// Get the key
	key, err := d.GetKey(token, id)
	if err != nil {
		return err
	}

	// Delete the key
	err = d.redisClient.Del(
		context.Background(),
		key,
	).Err()
	if err != nil {
		// Log the error
		if d.logger != nil {
			d.logger.DeleteTokenFromCacheFailed(err)
		}
	}
	return err
}
