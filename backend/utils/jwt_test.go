package utils_test

import (
	"fmt"
	"testing"
	"time"

	"backend/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"

	token, err := utils.GenerateJWT(userID, email)
	assert.NoError(t, err, "Failed to generate JWT")
	assert.NotEmpty(t, token, "Expected a non-empty token string")

	// Parse the token
	parsedToken, err := jwt.ParseWithClaims(token, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return utils.GetJWTKey(), nil
	})

	assert.NoError(t, err, "Error while parsing JWT token")
	assert.NotNil(t, parsedToken, "Parsed token should not be nil")
	assert.True(t, parsedToken.Valid, "Parsed token should be valid")

	// Extract claims
	claims, ok := parsedToken.Claims.(*utils.Claims)
	assert.True(t, ok, "Expected claims of type *utils.Claims")
	assert.Equal(t, userID, claims.UserID, "UserID does not match")
	assert.Equal(t, email, claims.Email, "Email does not match")
	assert.True(t, time.Now().Before(claims.ExpiresAt.Time), "Token should not be expired")
}

func TestGenerateJWTWithInvalidInputs(t *testing.T) {
	_, err := utils.GenerateJWT(0, "test@example.com")
	assert.Error(t, err, "Expected an error with invalid user ID")

	_, err = utils.GenerateJWT(1, "")
	assert.Error(t, err, "Expected an error with invalid email")
}

func TestExpiredJWT(t *testing.T) {
	// Create a token with an expiration time in the past
	expiredClaims := utils.Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, err := token.SignedString(utils.GetJWTKey())
	assert.NoError(t, err, "Error signing expired token")

	// Parse the expired token
	parsedToken, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTKey(), nil
	})

	// Assert that an error is returned when parsing an expired JWT token
	assert.Error(t, err, "Expected an error when parsing an expired JWT token")

	// The token object is not nil, but its Valid field should be false
	assert.False(t, parsedToken.Valid, "Expected token to be invalid (expired)")
}

func TestInvalidToken(t *testing.T) {
	_, err := utils.ParseJWT("invalid.token.here")
	assert.Error(t, err, "Expected an error when parsing an invalid JWT token")
}
