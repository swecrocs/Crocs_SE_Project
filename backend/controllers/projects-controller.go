package controllers

import (
	"net/http"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type ProjectCreationRequest struct {
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description"`
	RequiredSkills []string `json:"required_skills"`
	Visibility     string   `json:"visibility" binding:"oneof=private"`
	Status         string   `json:"status" binding:"oneof=open in-progress completed"`
}

type ProjectCreationResponse struct {
	Message string `json:"message" example:"Project successfully created"`
}

// CreateProject godoc
// @Summary      Create research project
// @Description  Creates a new research project and assigns the creator as an owner
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        request body ProjectCreationRequest true "Project attributes"
// @Success      201 {object} ProjectCreationResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /projects [post]
func CreateProject(c *gin.Context) {
	var request ProjectCreationRequest

	// validate request JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// extract user ID from the authenticated session
	userID := utils.InferUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
		return
	}

	// begin database transaction
	tx := database.DB.Begin()

	// create the project instance
	project := models.Project{
		Title:       request.Title,
		Description: request.Description,
		OwnerID:     userID,
		Visibility:  request.Visibility,
		Status:      request.Status,
	}
	project.SetRequiredSkills(request.RequiredSkills)

	// save to database within transaction
	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create project: " + err.Error()})
		return
	}

	// create a collaborator instance (owner)
	collaborator := models.Collaborator{
		ProjectID: project.ID,
		UserID:    userID,
		Role:      "owner",
	}

	// add the creator as the first collaborator
	if err := tx.Create(&collaborator).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to add project owner as collaborator: " + err.Error()})
		return
	}

	// commit the transaction
	if err := tx.Commit().Error; err != nil {
		// If commit fails, rollback and return error
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to commit transaction: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ProjectCreationResponse{Message: "Project successfully created"})
}
