package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	godatabasesredis "github.com/ralvarezdev/go-databases/redis"
	"time"
)

type (
	// TokenValidatorService struct
	TokenValidatorService struct {
		redisClient *redis.Client
	}
)

// NewTokenValidatorService creates a new token validator service
func NewTokenValidatorService(redisClient *redis.Client) (
	*TokenValidatorService,
	error,
) {
	// Check if the Redis client is nil
	if redisClient == nil {
		return nil, godatabasesredis.ErrNilClient
	}

	return &TokenValidatorService{redisClient: redisClient}, nil
}

// GetKey gets the JWT Identifier key
func (d *TokenValidatorService) GetKey(jti string) string {
	return godatabasesredis.GetKey(jti, JwtIdentifierPrefix)
}

// Set sets the token with the value and period
func (d *TokenValidatorService) Set(
	jti string,
	value interface{},
	expiresAt time.Time,
) error {
	// Get the key
	key := d.GetKey(jti)

	// Set the initial value
	if err := d.redisClient.Set(
		context.Background(),
		key,
		value,
		0,
	).Err(); err != nil {
		return err
	}

	// Set expiration time for the key as a UNIX timestamp
	return d.redisClient.ExpireAt(context.Background(), key, expiresAt).Err()
}

// Has checks if the token is valid
func (d *TokenValidatorService) Has(jti string) (bool, error) {
	// Get the key
	key := d.GetKey(jti)

	// Check the JWT Identifier
	_, err := d.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}

// Get gets the token
func (d *TokenValidatorService) Get(jti string) (interface{}, error) {
	// Get the key
	key := d.GetKey(jti)

	// Get the value
	value, err := d.redisClient.Get(
		context.Background(),
		key,
	).Result()
	if err != nil {
		return nil, err
	}
	return value, err
}

// Delete deletes the token
func (d *TokenValidatorService) Delete(jti string) error {
	// Get the key
	key := d.GetKey(jti)

	// Delete the key
	return d.redisClient.Del(
		context.Background(),
		key,
	).Err()
}
