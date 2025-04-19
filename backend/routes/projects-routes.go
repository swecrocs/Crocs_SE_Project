package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func ProjectsRoutes(router *gin.Engine) {
	projects := router.Group("/projects")
	{
		projects.GET("", controllers.ListProjects)
		projects.POST("", middleware.AuthRequired(), controllers.CreateProject)
		projects.GET("/:id", controllers.RetrieveProject)
		projects.POST("/:id/collaborators", middleware.AuthRequired(), controllers.InviteCollaborator)
		projects.GET("/invitations", middleware.AuthRequired(), controllers.GetProjectInvitations)
		projects.POST("/:id/collaborators/invitations/:invitationId/:action", middleware.AuthRequired(), controllers.RespondToProjectInvitation)
	}
}
