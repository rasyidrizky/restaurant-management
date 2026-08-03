package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims represents the JWT Claims
type Claims struct {
	UserID     uint `json:"user_id"`
	PositionID uint `json:"position_id"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID uint, positionID uint) (string, error) {
	// If secret is empty, fall back (for development only)
	secret := jwtSecret
	if len(secret) == 0 {
		secret = []byte("supersecret") 
	}

	claims := &Claims{
		UserID:     userID,
		PositionID: positionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ValidateJWT parses and validates a JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	secret := jwtSecret
	if len(secret) == 0 {
		secret = []byte("supersecret")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
