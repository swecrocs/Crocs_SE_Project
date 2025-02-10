package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	// Only login route for now
	// This handles user login
	router.POST("/auth/login", controllers.Login)
}
