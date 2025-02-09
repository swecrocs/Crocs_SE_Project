package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Authentication Routes
	router.POST("/auth/login", controllers.LoginUser)
}
