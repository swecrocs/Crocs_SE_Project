package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func ProjectsRoutes(router *gin.Engine) {
	projects := router.Group("/projects")
	{
		projects.POST("", middleware.AuthRequired(), controllers.CreateProject)
		projects.GET("/:id", controllers.RetrieveProject)
	}
}
