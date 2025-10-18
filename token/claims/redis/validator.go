package redis

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	godatabases "github.com/ralvarezdev/go-databases"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
)

type (
	// TokenValidatorService struct
	TokenValidatorService struct {
		redisClient *redis.Client
		logger      *slog.Logger
	}
)

// NewTokenValidatorService creates a new token validator service
//
// Parameters:
//
//   - redisClient: The Redis client
//   - logger: The logger (optional, can be nil)
//
// Returns:
//
//   - *TokenValidatorService: The token validator service
//   - error: An error if the Redis client is nil
func NewTokenValidatorService(
	redisClient *redis.Client,
	logger *slog.Logger,
) (
	*TokenValidatorService,
	error,
) {
	// Check if the Redis client is nil
	if redisClient == nil {
		return nil, godatabases.ErrNilConnection
	}

	if logger != nil {
		logger = logger.With(slog.String("component", "redis_token_validator"))
	}

	return &TokenValidatorService{redisClient, logger}, nil
}

// GetKey gets the JWT Identifier key
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - string: The key for the token
//   - error: An error if the token abbreviation fails
func (d *TokenValidatorService) GetKey(
	token gojwttoken.Token,
	id string,
) (string, error) {
	if d == nil {
		return "", gojwttokenclaims.ErrNilTokenValidator
	}

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
//
// Parameters:
//
//   - key: The key for the token
//   - isValid: The value to set (true if valid, false if revoked)
//   - expiresAt: The expiration time of the token
//
// Returns:
//
//   - error: An error if setting the token fails
func (d *TokenValidatorService) setWithFormattedKey(
	key string,
	isValid bool,
	expiresAt time.Time,
) error {
	if d == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Set the initial value
	if err := d.redisClient.Set(
		context.Background(),
		key,
		isValid,
		0,
	).Err(); err != nil {
		gojwttokenclaims.SetTokenFailed(err, d.logger)
		return err
	}

	// Set expiration time for the key as a UNIX timestamp
	err := d.redisClient.ExpireAt(context.Background(), key, expiresAt).Err()
	if err != nil {
		gojwttokenclaims.SetTokenFailed(err, d.logger)
	}
	return err
}

// Set sets the token with the value and expiration
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//   - isValid: The value to set (true if valid, false if revoked)
//   - expiresAt: The expiration time of the token
//
// Returns:
//
//   - error: An error if the token validator service is nil or if setting the token fails
func (d *TokenValidatorService) Set(
	token gojwttoken.Token,
	id string,
	isValid bool,
	expiresAt time.Time,
) error {
	if d == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := d.GetKey(token, id)
	if err != nil {
		return err
	}

	return d.setWithFormattedKey(key, isValid, expiresAt)
}

// Revoke revokes the token
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - error: An error if the token validator service is nil or if revoking the token fails
func (d *TokenValidatorService) Revoke(
	token gojwttoken.Token,
	id string,
) error {
	if d == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

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
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - bool: True if the token is valid, false if revoked
//   - error: An error if the token validator service is nil or if checking the token fails
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
		gojwttokenclaims.GetTokenFailed(err, d.logger)
		return false, err
	}

	// Parse the value
	parsedIsValue, err := strconv.ParseBool(isValid)
	if err != nil {
		return false, err
	}
	return parsedIsValue, nil
}
