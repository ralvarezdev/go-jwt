package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	godatabases "github.com/ralvarezdev/go-databases"
	gojwtcache "github.com/ralvarezdev/go-jwt/cache"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
	"strconv"
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
		return nil, godatabases.ErrNilConnection
	}

	return &TokenValidatorService{redisClient, logger}, nil
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
		JwtIdentifierSeparator,
		JwtIdentifierPrefix,
		tokenPrefix,
	), nil
}

// setWithFormattedKey sets the token with the value and expiration
func (d *TokenValidatorService) setWithFormattedKey(
	key string,
	isValid bool,
	expiresAt time.Time,
) error {
	// Set the initial value
	if err := d.redisClient.Set(
		context.Background(),
		key,
		isValid,
		0,
	).Err(); err != nil {
		// Log the error
		if d.logger != nil {
			d.logger.SetTokenToCacheFailed(err)
		}
		return err
	}

	// Set expiration time for the key as a UNIX timestamp
	err := d.redisClient.ExpireAt(context.Background(), key, expiresAt).Err()
	if err != nil {
		// Log the error
		if d.logger != nil {
			d.logger.SetTokenToCacheFailed(err)
		}
	}
	return err
}

// Set sets the token with the value and expiration
func (d *TokenValidatorService) Set(
	token gojwttoken.Token,
	id string,
	isValid bool,
	expiresAt time.Time,
) error {
	// Get the key
	key, err := d.GetKey(token, id)
	if err != nil {
		return err
	}

	return d.setWithFormattedKey(key, isValid, expiresAt)
}

// Revoke revokes the token
func (d *TokenValidatorService) Revoke(
	token gojwttoken.Token,
	id string,
) error {
	// Get the key
	key, err := d.GetKey(token, id)
	if err != nil {
		return err
	}

	// Get the current TTL of the key
	ttl, err := d.redisClient.TTL(context.Background(), key).Result()
	if err != nil {
		return err
	}

	// Update the value maintaining the TTL
	return d.setWithFormattedKey(key, false, time.Now().Add(ttl))
}

// IsValid checks if the token is valid
func (d *TokenValidatorService) IsValid(
	token gojwttoken.Token,
	id string,
) (bool, error) {
	// Get the key
	key, err := d.GetKey(token, id)
	if err != nil {
		return false, err
	}

	// Get the value
	isValid, err := d.redisClient.Get(
		context.Background(),
		key,
	).Result()
	if err != nil {
		// Log the error
		if d.logger != nil {
			d.logger.GetTokenFromCacheFailed(err)
		}
		return false, err
	}

	// Parse the value
	parsedIsValue, err := strconv.ParseBool(isValid)
	if err != nil {
		return false, err
	}
	return parsedIsValue, nil
}
