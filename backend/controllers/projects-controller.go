package controllers

import (
	"fmt"
	"net/http"
	"time"

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

type ProjectListResponse struct {
	Projects []ProjectRetrievalResponse `json:"projects"`
}

type MessageResponse struct {
	Message string `json:"message"`
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
			Status:         project.Status,
			OwnerID:        project.OwnerID,
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
	ID      uint   `json:"id"`
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

	c.JSON(http.StatusCreated, ProjectCreationResponse{
		Message: "Project successfully created",
		ID:      project.ID,
	})
}

type CollabInvitationRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=programmer editor"`
}

type CollabInvitationResponse struct {
	Message string `json:"message"`
}

// InviteCollaborator godoc
// @Summary      Invite a collaborator to a project
// @Description  Sends an invitation to a user to collaborate on a project
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        id path int true "Project ID"
// @Param        request body CollabInvitationRequest true "Invitation details"
// @Success      201 {object} CollabInvitationResponse
// @Failure      400 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /projects/{id}/collaborators [post]
func InviteCollaborator(c *gin.Context) {
	var request CollabInvitationRequest
	projectID := c.Param("id")

	// Validate request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get current user
	userID := utils.InferUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
		return
	}

	// Begin transaction
	tx := database.DB.Begin()

	// Verify project exists and user is owner
	var project models.Project
	if err := tx.First(&project, projectID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Project not found"})
		return
	}

	if project.OwnerID != userID {
		tx.Rollback()
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Only project owners can invite collaborators"})
		return
	}

	// Check if user is already a collaborator
	var existingCollaborator models.Collaborator
	if err := tx.Where("project_id = ? AND user_id = (SELECT id FROM users WHERE email = ?)",
		projectID, request.Email).First(&existingCollaborator).Error; err == nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, ErrorResponse{Error: "User is already a collaborator"})
		return
	}

	// Check for existing pending invitation
	var existingInvitation models.Invitation
	if err := tx.Where("project_id = ? AND email = ? AND status = ?",
		projectID, request.Email, models.InvitationStatusPending).First(&existingInvitation).Error; err == nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, ErrorResponse{Error: "An invitation is already pending"})
		return
	}

	// Create invitation
	invitation := models.Invitation{
		ProjectID: project.ID,
		InviterID: userID,
		Email:     request.Email,
		Role:      models.CollaboratorRole(request.Role),
		Status:    models.InvitationStatusPending,
	}

	if err := tx.Create(&invitation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create invitation"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to commit invitation"})
		return
	}

	c.JSON(http.StatusCreated, CollabInvitationResponse{
		Message: "Invitation sent successfully",
	})
}

type InvitationListResponse struct {
	Invitations []InvitationDetail `json:"invitations"`
}

type InvitationDetail struct {
	ID           uint      `json:"id"`
	ProjectID    uint      `json:"project_id"`
	ProjectTitle string    `json:"project_title"`
	InviterName  string    `json:"inviter_name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetProjectInvitations godoc
// @Summary      List pending invitations for the authenticated user
// @Description  Retrieves all pending invitations for the authenticated user
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Success      200 {object} InvitationListResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /projects/invitations [get]
func GetProjectInvitations(c *gin.Context) {
	userID := utils.InferUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication required"})
		return
	}

	// Get user's email
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch user"})
		return
	}

	// Get all pending invitations for this user's email
	var invitations []models.Invitation
	err := database.DB.Preload("Inviter").
		Preload("Project").
		Where("email = ? AND status = ?", user.Email, models.InvitationStatusPending).
		Find(&invitations).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch invitations"})
		return
	}

	response := make([]InvitationDetail, len(invitations))
	for i, inv := range invitations {
		response[i] = InvitationDetail{
			ID:           inv.ID,
			ProjectID:    inv.ProjectID,
			ProjectTitle: inv.Project.Title,
			InviterName:  inv.Inviter.Email,
			Email:        inv.Email,
			Role:         string(inv.Role),
			Status:       string(inv.Status),
			CreatedAt:    inv.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, InvitationListResponse{Invitations: response})
}

// RespondToProjectInvitation godoc
// @Summary      Accept or reject a project invitation
// @Description  Allows a user to accept or reject a collaboration invitation
// @Tags         Projects
// @Accept       json
// @Produce      json
// @Param        id path int true "Project ID"
// @Param        invitationId path int true "Invitation ID"
// @Param        action path string true "Action (accept/reject)"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      403 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /projects/{id}/collaborators/invitations/{invitationId}/{action} [post]
func RespondToProjectInvitation(c *gin.Context) {
	projectID := c.Param("id")
	invitationID := c.Param("invitationId")
	action := c.Param("action")
	userID := utils.InferUserID(c)

	if action != "accept" && action != "reject" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid action"})
		return
	}

	tx := database.DB.Begin()

	// Verify project exists
	var project models.Project
	if err := tx.First(&project, projectID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Project not found"})
		return
	}

	// Verify invitation exists and belongs to this project
	var invitation models.Invitation
	if err := tx.Where("id = ? AND project_id = ?", invitationID, projectID).First(&invitation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Invitation not found"})
		return
	}

	// Verify the invitation belongs to this user
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch user"})
		return
	}

	if invitation.Email != user.Email {
		tx.Rollback()
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Not authorized to respond to this invitation"})
		return
	}

	now := time.Now()
	invitation.ResponseDate = &now

	if action == "accept" {
		invitation.Status = models.InvitationStatusAccepted

		// Create collaborator entry
		collaborator := models.Collaborator{
			ProjectID: invitation.ProjectID,
			UserID:    userID,
			Role:      string(invitation.Role),
		}

		if err := tx.Create(&collaborator).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to add collaborator"})
			return
		}
	} else {
		invitation.Status = models.InvitationStatusRejected
	}

	if err := tx.Save(&invitation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update invitation"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, MessageResponse{Message: fmt.Sprintf("Invitation %s successfully", action+"ed")})
}
