package issuer

import (
	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
	"golang.org/x/crypto/ed25519"
)

type (
	// Ed25519Issuer handles JWT tokens issuing with ED25519 private key
	Ed25519Issuer struct {
		privateKey ed25519.PrivateKey
	}
)

// NewEd25519Issuer creates a new issuer by parsing the given path as an ED25519 private key
//
// Parameters:
//
//   - privateKey: The ED25519 private key in PEM format
//
// Returns:
//
//   - *Ed25519Issuer: The created issuer
//   - error: An error if the private key could not be parsed or is of an invalid type
func NewEd25519Issuer(privateKey []byte) (*Ed25519Issuer, error) {
	// Parse the private key
	key, err := jwt.ParseEdPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, gojwt.ErrUnableToParsePrivateKey
	}

	// Ensure the key is of type ED25519 private key
	ed25519Key, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, gojwt.ErrInvalidKeyType
	}

	return &Ed25519Issuer{
		privateKey: ed25519Key,
	}, nil
}

// IssueToken issues a new token for the given user with the given roles
//
// Parameters:
//
//   - claims: The claims to include in the token
//
// Returns:
//
//   - string: The issued token as a string
//   - error: An error if the token could not be created or signed
func (i Ed25519Issuer) IssueToken(claims jwt.Claims) (string, error) {
	// Create a new token with the claims
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, claims)

	// Sign and get the complete encoded token as a string using the private key
	rawToken, err := token.SignedString(i.privateKey)
	if err != nil {
		return "", err
	}

	return rawToken, nil
}
