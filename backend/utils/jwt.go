package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// TODO: switch to an environment variable in production
var jwtKey = []byte("your-secret-key")

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint, email string) (string, error) {
	// set expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// create claims
	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token with secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
