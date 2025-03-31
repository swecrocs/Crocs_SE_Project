package controllers

import (
	"net/http"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type ProjectRetrievalResponse struct {
	ID             uint     `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	RequiredSkills []string `json:"required_skills"`
	Visibility     string   `json:"visibility"`
	Status         string   `json:"status"`
	OwnerID        uint     `json:"owner_id"`
}

// GetProject godoc
// @Summary      Get research project details
// @Description  Retrieves details of a specific research project by ID
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        id path int true "Project ID"
// @Success      200 {object} ProjectRetrievalResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /projects/{id} [get]
func RetrieveProject(c *gin.Context) {
	// Get project ID from URL parameter
	projectID := c.Param("id")

	// Find project in database
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch project: " + err.Error()})
		return
	}

	// Convert to response format
	response := ProjectRetrievalResponse{
		ID:             project.ID,
		Title:          project.Title,
		Description:    project.Description,
		RequiredSkills: project.GetRequiredSkills(),
		Visibility:     project.Visibility,
		Status:         project.Status,
		OwnerID:        project.OwnerID,
	}

	c.JSON(http.StatusOK, response)
}

type ProjectListResponse struct {
    Projects []ProjectRetrievalResponse `json:"projects"`
}

// ListProjects godoc
// @Summary      List all research projects
// @Description  Retrieves a list of all research projects
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Success      200 {object} ProjectListResponse
// @Failure      500 {object} ErrorResponse
// @Router       /projects [get]
func ListProjects(c *gin.Context) {
    var projects []models.Project
    
    // Fetch all projects with their owners
    if err := database.DB.Find(&projects).Error; err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch projects"})
        return
    }

    // Convert to response format
    response := make([]ProjectRetrievalResponse, len(projects))
    for i, project := range projects {
        response[i] = ProjectRetrievalResponse{
            ID:             project.ID,
            Title:          project.Title,
            Description:    project.Description,
            RequiredSkills: project.GetRequiredSkills(),
            Visibility:     project.Visibility,
            Status:        project.Status,
            OwnerID:       project.OwnerID,
        }
    }

    c.JSON(http.StatusOK, ProjectListResponse{Projects: response})
}

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
