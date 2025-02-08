package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.RegisterUser)
	}
}
