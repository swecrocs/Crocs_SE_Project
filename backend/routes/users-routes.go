package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.PUT("/:id", controllers.EditUserProfile)
	}
}
