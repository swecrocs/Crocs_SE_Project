package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// Register routes
func RegisterRoutes(router *gin.Engine) {
	// Auth routes
	router.POST("/auth/login", controllers.Login)

	// User routes with JWT authentication
	router.PUT("/users/:id/profile", middleware.JWTMiddleware(), controllers.EditUserProfile)
}
