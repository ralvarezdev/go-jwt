package claims

import (
	"github.com/golang-jwt/jwt/v5"
	gojwt "github.com/ralvarezdev/go-jwt"
	gojwtcache "github.com/ralvarezdev/go-jwt/cache"
	gojwttoken "github.com/ralvarezdev/go-jwt/token"
)

type (
	// DefaultClaimsValidator is the default implementation of the ClaimsValidator interface
	DefaultClaimsValidator struct {
		tokenValidator TokenValidator
	}
)

// NewDefaultClaimsValidator creates a new default claims validator
//
// Parameters:
//
//   - tokenValidator: the token validator
//
// Returns:
//
//   - *DefaultValidator: the default validator
//   - error: if there was an error creating the default validator
func NewDefaultClaimsValidator(
	tokenValidator TokenValidator,
) (*DefaultClaimsValidator, error) {
	// Check if the token validator is nil
	if tokenValidator == nil {
		return nil, gojwtcache.ErrNilTokenValidator
	}

	return &DefaultClaimsValidator{
		tokenValidator,
	}, nil
}

// ValidateClaims validates the claims
//
// Parameters:
//
//   - claims: the claims to validate
//   - token: the token type (access or refresh)
//
// Returns:
//
//   - bool: true if the claims are valid, false otherwise
//   - error: if there was an error validating the claims
func (d DefaultClaimsValidator) ValidateClaims(
	claims jwt.MapClaims,
	token gojwttoken.Token,
) (bool, error) {
	// Check if the claims are nil
	if claims == nil {
		return false, ErrNilClaims
	}

	// Get the JWT Identifier
	jti, ok := claims[gojwt.IDClaim]
	if !ok {
		return false, ErrIDClaimNotFound
	}

	// Check if the JWT Identifier is a string
	jtiStr, ok := jti.(string)
	if !ok {
		return false, ErrInvalidIDClaim
	}

	// Check if the token is valid
	isValid, err := d.tokenValidator.IsValid(token, jtiStr)
	if err != nil {
		return false, err
	}
	return isValid, nil
}
