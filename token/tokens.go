package token

// Token represents a token type
type Token string

var (
	// RefreshToken represents a refresh token
	RefreshToken Token = "refresh_token"

	// AccessToken represents an access token
	AccessToken Token = "access_token"
)

// String returns the string representation of the token
func (t Token) String() string {
	return string(t)
}

// Abbreviation returns the abbreviation of the token
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
