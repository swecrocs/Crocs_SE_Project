package middleware

import (
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthResponse represents standardized auth error responses
type AuthResponse struct {
	Error string `json:"error"`
}

// AuthRequired validates JWT tokens and sets user ID in context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, AuthResponse{Error: "Authorization header is required"})
			c.Abort()
			return
		}

		// Check for Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, AuthResponse{Error: "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Parse and validate the token
		tokenString := parts[1]
		claims, err := utils.ParseJWT(tokenString)
		
		if err != nil {
			c.JSON(http.StatusUnauthorized, AuthResponse{Error: "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(utils.UserIDKey, claims.UserID)
		c.Set(utils.UserEmailKey, claims.Email)
		
		c.Next()
	}
}