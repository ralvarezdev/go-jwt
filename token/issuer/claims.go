package issuer

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// GetExpirationTime returns the expiration time for the given duration
func GetExpirationTime(
	issuedTime time.Time,
	duration time.Duration,
) time.Time {
	return issuedTime.Add(duration)
}

// GenerateClaims generates a new claims object
func GenerateClaims(
	issuedAt time.Time,
	expirationAt time.Time,
	additionalClaims map[string]interface{},
) *jwt.MapClaims {
	claims := jwt.MapClaims{
		"exp": expirationAt.Unix(),
		"iat": issuedAt.Unix(),
	}

	// Add the additional claims
	for key, value := range additionalClaims {
		claims[key] = value
	}

	return &claims
}
