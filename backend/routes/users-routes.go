package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.GET("/:id", controllers.RetrieveUserProfile)
		users.PUT("/:id", controllers.EditUserProfile)
	}
}
