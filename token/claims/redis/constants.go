package redis

import (
	gostringsseparator "github.com/ralvarezdev/go-strings/separator"
)

var (
	// ParentRefreshTokenIDPrefix is the prefix of the Parent Refresh Token ID key
	ParentRefreshTokenIDPrefix = "prt"

	// KeySeparator is the separator for the Redis keys
	KeySeparator = gostringsseparator.Dots
)
