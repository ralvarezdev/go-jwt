package redis

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	godatabases "github.com/ralvarezdev/go-databases"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
	"github.com/redis/go-redis/v9"
)

type (
	// TokenValidator struct
	TokenValidator struct {
		redisClient *redis.Client
		logger      *slog.Logger
	}
)

// NewTokenValidator creates a new token validator
//
// Parameters:
//
//   - redisClient: The Redis client
//   - logger: The logger (optional, can be nil)
//
// Returns:
//
//   - *TokenValidator: The token validator
//   - error: An error if the Redis client is nil
func NewTokenValidator(
	redisClient *redis.Client,
	logger *slog.Logger,
) (
	*TokenValidator,
	error,
) {
	// Check if the Redis client is nil
	if redisClient == nil {
		return nil, godatabases.ErrNilConnection
	}

	if logger != nil {
		logger = logger.With(slog.String("component", "redis_token_validator"))
	}

	return &TokenValidator{redisClient, logger}, nil
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
func (t *TokenValidator) GetKey(
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

	return gostringsadd.Prefixes(
		id,
		KeySeparator,
		tokenPrefix,
	), nil
}

// GetParentRefreshTokenKey gets the parent refresh token key
//
// Parameters:
//
//   - id: The ID associated with the refresh token
//
// Returns:
//
//   - string: The key for the parent refresh token
//   - error: An error if the token validator is nil
func (t *TokenValidator) GetParentRefreshTokenKey(
	id string,
) (string, error) {
	if t == nil {
		return "", gojwttokenclaims.ErrNilTokenValidator
	}

	return gostringsadd.Prefixes(
		id,
		KeySeparator,
		ParentRefreshTokenIDPrefix,
	), nil
}

// setKey sets the token with the value and expiration
//
// Parameters:
//
//   - ctx: The context
//   - key: The key for the token
//   - isValid: The value to set (true if valid, false if revoked)
//   - expiresAt: The expiration time of the token
//
// Returns:
//
//   - error: An error if setting the token fails
func (t *TokenValidator) setKey(
	ctx context.Context,
	key string,
	isValid bool,
	expiresAt time.Time,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Set the initial value
	if err := t.redisClient.Set(
		ctx,
		key,
		isValid,
		0,
	).Err(); err != nil {
		gojwttokenclaims.SetTokenFailed(err, t.logger)
		return err
	}

	// Set expiration time for the key as a UNIX timestamp
	err := t.redisClient.ExpireAt(ctx, key, expiresAt).Err()
	if err != nil {
		gojwttokenclaims.SetTokenFailed(err, t.logger)
	}
	return err
}

// AddRefreshToken adds a refresh token
//
// Parameters:
//
//   - id: The ID associated with the token
//   - expiresAt: The expiration time of the token
//
// Returns:
//
//   - error: An error if the token validator is nil or if adding the refresh token fails
func (t *TokenValidator) AddRefreshToken(
	ctx context.Context,
	id string,
	expiresAt time.Time,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := t.GetKey(gojwttoken.RefreshToken, id)
	if err != nil {
		return err
	}

	return t.setKey(ctx, key, true, expiresAt)
}

// AddAccessToken adds an access token
//
// Parameters:
//
//   - ctx: The context
//   - id: The ID associated with the token
//   - parentRefreshTokenID: The parent refresh token ID
//   - expiresAt: The expiration time of the token
//
// Returns:
//
//   - error: An error if the token validator is nil or if adding the access token fails
func (t *TokenValidator) AddAccessToken(
	ctx context.Context,
	id string,
	parentRefreshTokenID string,
	expiresAt time.Time,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := t.GetKey(gojwttoken.AccessToken, id)
	if err != nil {
		return err
	}

	// Set the parent refresh token ID key
	parentRefreshTokenKey, parentKeyErr := t.GetParentRefreshTokenKey(parentRefreshTokenID)
	if parentKeyErr != nil {
		return parentKeyErr
	}

	// Set the parent refresh token ID with its access token ID
	if setErr := t.redisClient.Set(
		ctx,
		parentRefreshTokenKey,
		id,
		time.Until(expiresAt),
	).Err(); setErr != nil {
		gojwttokenclaims.SetTokenFailed(setErr, t.logger)
		return setErr
	}

	return t.setKey(ctx, key, true, expiresAt)
}

// RevokeToken revokes the token
//
// Parameters:
//
//   - ctx: The context
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - error: An error if the token validator is nil or if revoking the token fails
func (t *TokenValidator) RevokeToken(
	ctx context.Context,
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

	// Get the current TTL of the key
	ttl, err := t.redisClient.TTL(ctx, key).Result()
	if err != nil {
		return err
	}

	// Update the value maintaining the TTL
	if err = t.setKey(
		ctx,
		key,
		false,
		time.Now().Add(ttl),
	); err != nil {
		gojwttokenclaims.RevokeTokenFailed(err, t.logger)
		return err
	}

	// Check if the token is a refresh token to revoke its associated access token
	if token == gojwttoken.AccessToken {
		return nil
	}

	// Get the parent refresh token key
	parentRefreshTokenKey, err := t.GetParentRefreshTokenKey(id)
	if err != nil {
		return err
	}

	// Get the associated access token ID
	accessTokenID, err := t.redisClient.Get(
		ctx,
		parentRefreshTokenKey,
	).Result()
	if err != nil {
		return err
	}

	// Revoke the associated access token
	accessTokenKey, err := t.GetKey(gojwttoken.AccessToken, accessTokenID)
	if err != nil {
		return err
	}

	// Get the current TTL of the access token key
	accessTokenTTL, err := t.redisClient.TTL(
		ctx,
		accessTokenKey,
	).Result()
	if err != nil {
		return err
	}

	// Update the value maintaining the TTL
	if err = t.setKey(
		ctx,
		accessTokenKey,
		false,
		time.Now().Add(accessTokenTTL),
	); err != nil {
		gojwttokenclaims.RevokeTokenFailed(err, t.logger)
		return err
	}
	return nil
}

// IsTokenValid checks if the token is valid
//
// Parameters:
//
//   - ctx: The context
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - bool: True if the token is valid, false if revoked
//   - error: An error if the token validator is nil or if checking the token fails
func (t *TokenValidator) IsTokenValid(
	ctx context.Context,
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

	// Get the value
	isValid, err := t.redisClient.Get(
		ctx,
		key,
	).Result()
	if err != nil {
		// Check if the error is a redis.Nil error (key does not exist)
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		gojwttokenclaims.GetTokenFailed(err, t.logger)
		return false, err
	}

	// Parse the value
	parsedIsValue, err := strconv.ParseBool(isValid)
	if err != nil {
		return false, err
	}
	return parsedIsValue, nil
}
