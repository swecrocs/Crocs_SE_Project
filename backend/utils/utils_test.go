package utils

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	// Sample input data
	userID := uint(1)
	email := "test@example.com"

	// Call the GenerateJWT function
	token, err := GenerateJWT(userID, email)

	// Assert that there is no error
	assert.NoError(t, err, "Expected no error while generating JWT token")

	// Assert that the token string is not empty
	assert.NotEmpty(t, token, "Expected a non-empty token string")

	// Decode the token to check its structure (optional)
	// This will split the token into 3 parts (header, payload, signature) by "."
	parts := strings.Split(token, ".")
	assert.Len(t, parts, 3, "JWT token should have 3 parts")

	// Decode the token and verify the claims
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Return the JWT key for signing validation
		return jwtKey, nil
	})

	// Assert that parsing the token was successful
	assert.NoError(t, err, "Error while parsing JWT token")

	// Assert that the token is valid and contains the correct claims
	claims, ok := parsedToken.Claims.(*Claims)
	assert.True(t, ok, "The token should contain valid claims")

	// Check if the claims are correctly set
	assert.Equal(t, userID, claims.UserID, "The UserID in the token is incorrect")
	assert.Equal(t, email, claims.Email, "The Email in the token is incorrect")

	// Check if the expiration time is in the future
	assert.True(t, time.Now().Before(time.Unix(claims.ExpiresAt, 0)), "The token should not have expired")
}
