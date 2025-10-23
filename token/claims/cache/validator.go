package cache

import (
	"context"
	"log/slog"
	"time"

	gocache "github.com/ralvarezdev/go-cache"
	gocachetimed "github.com/ralvarezdev/go-cache/timed"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwttokenclaims "github.com/ralvarezdev/go-jwt/token/claims"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
)

type (
	// TokenValidator struct
	TokenValidator struct {
		logger *slog.Logger
		cache  gocachetimed.TimedCache
	}
)

// NewTokenValidator creates a new token validator
//
// Parameters:
//
//   - logger: The logger (optional, can be nil)
//
// Returns:
//
//   - *TokenValidator: The token validator
func NewTokenValidator(logger *slog.Logger) *TokenValidator {
	if logger != nil {
		logger = logger.With(
			slog.String(
				"component",
				"cache_token_validator",
			),
		)
	}
	return &TokenValidator{
		cache:  gocachetimed.NewDefaultTimedCache(),
		logger: logger,
	}
}

// GetTokenKey gets the JWT Identifier key
//
// Parameters:
//
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - string: The key for the cache
//   - error: An error if the token validator is nil or if the token abbreviation fails
func (t *TokenValidator) GetTokenKey(
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

	return gostringsadd.Prefixes(tokenPrefix, KeySeparator, id), nil
}

// GetParentRefreshTokenKey gets the parent
//
// Parameters:
//
//   - id: The ID of the refresh token
//   - parentTokenPrefix: The parent token prefix
//
// Returns:
//
//   - string: The key for the cache
//   - error: An error if the token validator is nil
func (t *TokenValidator) GetParentRefreshTokenKey(
	id string,
) (string, error) {
	if t == nil {
		return "", gojwttokenclaims.ErrNilTokenValidator
	}

	return gostringsadd.Prefixes(
		ParentRefreshTokenIDPrefix,
		KeySeparator,
		id,
	), nil
}

// AddRefreshToken sets a token in the cache
//
// Parameters:
//
//   - ctx: The context (not used, but kept for interface consistency)
//   - token: The token
//   - id: The ID associated with the token
//   - expiresAt: The expiration time of the token
//
// Returns:
//
//   - error: An error if the token validator is nil or if setting the token in the cache fails
func (t *TokenValidator) AddRefreshToken(
	ctx context.Context,
	id string,
	expiresAt time.Time,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := t.GetTokenKey(gojwttoken.RefreshToken, id)
	if err != nil {
		return err
	}

	// Set the token in the cache
	err = t.cache.Set(key, gocachetimed.NewTimedItem(true, expiresAt))
	if err != nil {
		gojwttokenclaims.SetTokenFailed(err, t.logger)
	}
	return err
}

// AddAccessToken sets a token in the cache
//
// Parameters:
//
//   - ctx: The context (not used, but kept for interface consistency)
//   - id: The ID associated with the token
//   - parentRefreshTokenID: The parent refresh token ID
//   - expiresAt: The expiration time of the token
//
// Returns:
//
//   - error: An error if the token validator is nil or if setting the token in the cache fails
func (t *TokenValidator) AddAccessToken(
	ctx context.Context,
	id string,
	parentRefreshTokenID string,
	expiresAt time.Time,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Check if the parent refresh token has already been set
	refreshTokenKey, err := t.GetTokenKey(
		gojwttoken.RefreshToken,
		parentRefreshTokenID,
	)
	if err != nil {
		return err
	}

	// Get the parent refresh token from the cache
	value, found := t.cache.Get(refreshTokenKey)
	if !found {
		return ErrParentRefreshTokenNotFound
	}

	// Parse the value to check if it's valid
	refreshTokenItem, ok := value.(*gocachetimed.TimedItem)
	if !ok {
		return ErrInvalidParentRefreshTokenItem
	}

	// Check if the parent refresh token is still valid
	if refreshTokenItem.HasExpired() {
		return nil
	}

	// Get the key
	key, err := t.GetTokenKey(gojwttoken.AccessToken, id)
	if err != nil {
		return err
	}

	// Set the token in the cache
	err = t.cache.Set(key, gocachetimed.NewTimedItem(true, expiresAt))
	if err != nil {
		gojwttokenclaims.SetTokenFailed(err, t.logger)
		return err
	}

	// Also set the parent refresh token key to point to this access token
	parentRefreshTokenKey, err := t.GetParentRefreshTokenKey(
		parentRefreshTokenID,
	)
	if err != nil {
		return err
	}

	// Set the parent refresh token in the cache
	if err = t.cache.Set(
		parentRefreshTokenKey,
		gocachetimed.NewTimedItem(id, expiresAt),
	); err != nil {
		gojwttokenclaims.SetTokenFailed(err, t.logger)
	}
	return err
}

// RevokeToken revokes a token in the cache
//
// Parameters:
//
//   - ctx: The context (not used, but kept for interface consistency)
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - error: An error if the token validator is nil or if revoking the token in the cache fails
func (t *TokenValidator) RevokeToken(
	ctx context.Context,
	token gojwttoken.Token,
	id string,
) error {
	if t == nil {
		return gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := t.GetTokenKey(token, id)
	if err != nil {
		return err
	}

	// Get the item from the cache to check if it exists
	value, found := t.cache.Get(key)
	if !found {
		return gocache.ErrItemNotFound
	}

	// Check if the item is a TimedItem
	item, ok := value.(*gocachetimed.TimedItem)
	if !ok {
		return ErrInvalidTokenItem
	}

	// Revoke the token in the cache
	item.SetValue(false)

	// Also, revoke the access token if it's a refresh token
	if token != gojwttoken.RefreshToken {
		return nil
	}

	// Get the parent refresh token key
	parentKey, err := t.GetParentRefreshTokenKey(
		id,
	)
	if err != nil {
		return err
	}

	// Get the access token ID from the parent refresh token key
	value, found = t.cache.Get(parentKey)
	if !found {
		return nil
	}

	// Parse the value to get the access token ID
	parentRefreshTokenItem, ok := value.(*gocachetimed.TimedItem)
	if !ok {
		return ErrInvalidParentRefreshTokenItem
	}

	// Get the access token ID from the item
	accessTokenID, ok := parentRefreshTokenItem.GetValue().(string)
	if !ok {
		return ErrInvalidParentRefreshTokenItem
	}

	// Revoke the access token in the cache
	return t.RevokeToken(
		ctx,
		gojwttoken.AccessToken,
		accessTokenID,
	)
}

// IsTokenValid checks if a token is valid in the cache
//
// Parameters:
//
//   - ctx: The context (not used, but kept for interface consistency)
//   - token: The token
//   - id: The ID associated with the token
//
// Returns:
//
//   - bool: Whether the token is valid
//   - error: An error if the token validator is nil or if checking the token in the cache fails
func (t *TokenValidator) IsTokenValid(
	ctx context.Context,
	token gojwttoken.Token,
	id string,
) (bool, error) {
	if t == nil {
		return false, gojwttokenclaims.ErrNilTokenValidator
	}

	// Get the key
	key, err := t.GetTokenKey(token, id)
	if err != nil {
		return false, err
	}

	// Get the token from the cache
	value, found := t.cache.Get(key)
	if !found {
		return false, gocache.ErrItemNotFound
	}

	// Check if the value is a TimedItem
	timedItem, ok := value.(*gocachetimed.TimedItem)
	if !ok {
		return false, ErrInvalidTokenItem
	}

	// Check if the item has expired
	if timedItem.HasExpired() {
		return false, nil
	}

	// Return the validity of the token
	isValid, ok := timedItem.GetValue().(bool)
	if !ok {
		return false, ErrInvalidTokenItem
	}

	return isValid, nil
}
