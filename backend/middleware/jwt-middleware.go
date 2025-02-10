package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// JWTMiddleware to protect routes
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse and validate the JWT token
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for JWT signing
			return []byte("your-secret-key"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid or missing token"})
			c.Abort()
			return
		}

		// Continue to the handler if token is valid
		c.Next()
	}
}
