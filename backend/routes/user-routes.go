package routes

import (
	"github.com/gin-gonic/gin"
	"backend/controllers"
)

func UserRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.RegisterUser)
	}
}