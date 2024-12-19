package token

// Token represents a token type
type Token string

const (
	// RefreshToken represents a refresh token
	RefreshToken Token = "refresh_token"

	// AccessToken represents an access token
	AccessToken Token = "access_token"
)

// String returns the string representation of the token
func (t Token) String() string {
	return string(t)
}
