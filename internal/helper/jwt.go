package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT generates a JWT token for a given user ID and secret key.
func GenerateJWT(userID int64, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
