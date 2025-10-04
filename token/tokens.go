package token

type (
	// Token represents a token type
	Token string
)

var (
	// RefreshToken represents a refresh token
	RefreshToken Token = "refresh_token"

	// AccessToken represents an access token
	AccessToken Token = "access_token"
)

// String returns the string representation of the token
//
// Returns:
//
//   - string: The string representation of the token
func (t Token) String() string {
	return string(t)
}

// Abbreviation returns the abbreviation of the token
//
// Returns:
//
//   - string: The abbreviation of the token
//   - error: An error if the token type is unexpected
func (t Token) Abbreviation() (string, error) {
	switch t {
	case RefreshToken:
		return "RT", nil
	case AccessToken:
		return "AT", nil
	default:
		return "", ErrUnexpectedTokenType
	}
}
