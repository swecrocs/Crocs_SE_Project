package middleware

import (
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SameUserOnly ensures that the authenticated user can only access/modify their own resources
func SameUserOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated user ID from context
		authUserID := utils.InferUserID(c)
		if authUserID == 0 {
			c.JSON(http.StatusUnauthorized, AuthResponse{Error: "Authentication required"})
			c.Abort()
			return
		}

		// Get requested user ID from URL
		requestedUserID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, AuthResponse{Error: "Invalid user ID"})
			c.Abort()
			return
		}

		// Check if authenticated user is requesting their own resource
		if authUserID != uint(requestedUserID) {
			c.JSON(http.StatusForbidden, AuthResponse{Error: "You can only modify your own profile"})
			c.Abort()
			return
		}

		c.Next()
	}
}