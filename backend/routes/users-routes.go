package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.GET("/:id/profile", controllers.RetrieveUserProfile)
		users.PUT("/:id/profile", middleware.AuthRequired(), middleware.SameUserOnly(), controllers.EditUserProfile)
	}
}
