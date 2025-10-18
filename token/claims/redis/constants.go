package redis

import (
	gostringsseparator "github.com/ralvarezdev/go-strings/separator"
)

var (
	// JwtIdentifierPrefix is the prefix of the JWT Identifier key
	JwtIdentifierPrefix = "jti"

	// JwtIdentifierSeparator is the separator of the JWT identifier prefixes
	JwtIdentifierSeparator = gostringsseparator.Dots
)
