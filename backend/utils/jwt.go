package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key_here") // Store in .env for better security

// Claims structure for JWT
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT generates JWT token
func GenerateJWT(email string) (string, error) {
	// Create the JWT claims, which includes the username and expiration time
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
