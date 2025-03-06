package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT secret key from environment variable with fallback
var jwtKey = []byte(getJWTKey())

// Claims structure for JWT
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GetJWTKey returns the JWT signing key
func GetJWTKey() []byte {
	return jwtKey
}

// getJWTKey retrieves JWT key from environment or uses fallback
func getJWTKey() string {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		// Fallback for development - NEVER use this in production
		return "dev-temporary-key-replace-in-production"
	}
	return key
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID uint, email string) (string, error) {
	if userID == 0 || email == "" {
		return "", fmt.Errorf("invalid input: userID and email are required")
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1 day expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseJWT parses and validates a JWT token string
func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
