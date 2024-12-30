package auth

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	godatabasesredis "github.com/ralvarezdev/go-databases/redis"
	"strconv"
	"time"
)

type (
	// TokenValidator interface
	TokenValidator interface {
		AddToken(jwtId string, period time.Duration) error
		RevokeToken(jwtId string) error
		IsTokenValid(jwtId string) (bool, error)
	}

	// DefaultTokenValidator struct
	DefaultTokenValidator struct {
		redisClient *redis.Client
	}
)

// NewDefaultTokenValidator creates a new token validator
func NewDefaultTokenValidator(redisClient *redis.Client) (
	*DefaultTokenValidator,
	error,
) {
	// Check if the Redis client is nil
	if redisClient == nil {
		return nil, godatabasesredis.ErrNilClient
	}

	return &DefaultTokenValidator{redisClient: redisClient}, nil
}

// GetKey gets the JWT Identifier key
func (d *DefaultTokenValidator) GetKey(jwtId string) string {
	return godatabasesredis.GetKey(jwtId, JwtIdentifierPrefix)
}

// AddToken adds the JWT Identifier with the given expiration period
func (d *DefaultTokenValidator) AddToken(
	jwtId string,
	period time.Duration,
) error {
	key := d.GetKey(jwtId)

	// Set the initial value
	_, err := d.redisClient.Set(
		context.Background(),
		key,
		true,
		period,
	).Result()
	return err
}

// RevokeToken revokes the JWT Identifier by setting the value to false
func (d *DefaultTokenValidator) RevokeToken(jwtId string) error {
	key := d.GetKey(jwtId)

	// Set the token as invalid by setting the value to false
	_, err := d.redisClient.Set(
		context.Background(),
		key,
		false,
		redis.KeepTTL,
	).Result()
	return err
}

// IsTokenValid checks if the token is valid
func (d *DefaultTokenValidator) IsTokenValid(jwtId string) (bool, error) {
	key := d.GetKey(jwtId)

	// Check the JWT Identifier
	value, err := d.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	// Parse value
	isValid, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf(ErrFailedToParseBoolValue, err)
	}

	return isValid, nil
}
