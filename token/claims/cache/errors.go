package cache

import (
	"errors"
)

var (
	ErrParentRefreshTokenNotFound    = errors.New("parent refresh token not found")
	ErrInvalidParentRefreshTokenItem = errors.New("invalid parent refresh token item")
	ErrInvalidTokenItem              = errors.New("invalid token item")
)
