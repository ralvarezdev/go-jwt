package redis

import (
	gostringsadd "github.com/ralvarezdev/go-strings/add"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

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
func GetKey(
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
func GetParentRefreshTokenKey(
	id string,
) string {
	return gostringsadd.Prefixes(
		id,
		KeySeparator,
		ParentRefreshTokenIDPrefix,
	)
}