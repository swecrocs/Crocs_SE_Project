package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.GET("/:id/profile", controllers.RetrieveUserProfile)
		users.PUT("/:id/profile", controllers.EditUserProfile)
	}
}
