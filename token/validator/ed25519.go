package validator

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	goflagmode "github.com/ralvarezdev/go-flags/mode"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
	gojwtclaims "github.com/ralvarezdev/go-jwt/token/claims"
	"golang.org/x/crypto/ed25519"
)

type (
	// Ed25519Validator handles parsing and validation of JWT tokens with ED25519 public key
	Ed25519Validator struct {
		publicKey       ed25519.PublicKey
		claimsValidator gojwtclaims.Validator
		mode            *goflagmode.Flag
	}
)

// NewEd25519Validator returns a new validator by parsing the given file path as an ED25519 public key
//
// Parameters:
//
//   - publicKey: The ED25519 public key in PEM format
//   - claimsValidator: The token claims validator
//   - mode: The mode flag to determine if debug mode is enabled
//
// Returns:
//
//   - *Ed25519Validator: The ED25519 validator
//   - error: An error if the public key cannot be parsed or if any parameter is nil
func NewEd25519Validator(
	publicKey []byte,
	claimsValidator gojwtclaims.Validator,
	mode *goflagmode.Flag,
) (*Ed25519Validator, error) {
	// Check if either the token validator or the mode flag is nil
	if claimsValidator == nil {
		return nil, ErrNilClaimsValidator
	}
	if mode == nil {
		return nil, goflagmode.ErrNilModeFlag
	}

	// Parse the public key
	key, err := jwt.ParseEdPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, gojwt.ErrUnableToParsePublicKey
	}

	// Ensure the key is of type ED25519 public key
	ed25519Key, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, gojwt.ErrInvalidKeyType
	}

	return &Ed25519Validator{
		publicKey:       ed25519Key,
		claimsValidator: claimsValidator,
		mode:            mode,
	}, nil
}

// GetToken parses the given JWT raw token
//
// Parameters:
//
//   - rawToken: The raw JWT token string
//
// Returns:
//
//   - *jwt.Token: The parsed JWT token
//   - error: An error if the token is invalid or if parsing fails
func (d Ed25519Validator) GetToken(rawToken string) (*jwt.Token, error) {
	// Parse JWT and verify signature
	token, err := jwt.Parse(
		rawToken,
		func(rawToken *jwt.Token) (interface{}, error) {
			// Check to see if the token uses the expected signing method
			if _, ok := rawToken.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return d.publicKey, nil
		},
	)
	if err != nil {
		// Check if the mode is debug
		if d.mode != nil && d.mode.IsDebug() {
			return nil, err
		}

		switch {
		case errors.Is(err, ErrUnexpectedSigningMethod):
		case errors.Is(err, jwt.ErrSignatureInvalid):
		case errors.Is(err, jwt.ErrTokenExpired):
		case errors.Is(err, jwt.ErrTokenNotValidYet):
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, err
		default:
			return nil, ErrInvalidToken
		}
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return token, nil
}

// GetClaims parses and validates the given JWT raw token
//
// Parameters:
//
//   - rawToken: The raw JWT token string
//
// Returns:
//
//   - jwt.MapClaims: The token claims
//   - error: An error if the token is invalid, if parsing fails, or if the claims are of an unexpected type
func (d Ed25519Validator) GetClaims(rawToken string) (
	jwt.MapClaims, error,
) {
	// Get the token
	token, err := d.GetToken(rawToken)
	if err != nil {
		return nil, err
	}

	// Get token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// ValidateClaims validates the given token claims based on the given token type and returns the claims if valid
//
// Parameters:
//
//   - rawToken: The raw JWT token string
//   - token: The token type
//
// Returns:
//
//   - jwt.MapClaims: The token claims if valid
//   - error: An error if the token is invalid, if parsing fails, or if the claims are invalid
func (d Ed25519Validator) ValidateClaims(
	rawToken string,
	token gojwttoken.Token,
) (jwt.MapClaims, error) {
	// Get the claims
	claims, err := d.GetClaims(rawToken)
	if err != nil {
		return nil, err
	}

	// Check if the token claims are valid
	areValid, err := d.claimsValidator.ValidateClaims(claims, token)
	if err != nil {
		return nil, err
	}
	if !areValid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
