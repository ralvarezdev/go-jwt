package cache

import (
	gostringsseparator "github.com/ralvarezdev/go-strings/separator"
)

var (
	// KeySeparator is the separator for the cache keys
	KeySeparator = gostringsseparator.Dots

	// ParentRefreshTokenIDPrefix is the prefix of the Parent Refresh Token ID key
	ParentRefreshTokenIDPrefix = "PRT"
)
