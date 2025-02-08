package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileEditRequest struct {
	FullName    string `json:"full_name"`
	Bio         string `json:"bio"`
	Affiliation string `json:"affiliation"`
}

type ProfileEditResponse struct {
	Message string `json:"message"`
}

// EditUserProfile godoc
// @Summary      Edit user profile
// @Description  Update an existing user profile with new information.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        request body ProfileEditRequest true "Profile information"
// @Success      200 {object} ProfileEditResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{id}/profile [put]
func EditUserProfile(c *gin.Context) {
	userID := c.Param("id")

	// find the profile linked to this user
	var profile models.UserProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Profile not found"})
		return
	}

	// bind request JSON to ProfileEditRequest struct
	var request ProfileEditRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	// update profile fields
	profile.FullName = request.FullName
	profile.Bio = request.Bio
	profile.Affiliation = request.Affiliation

	// save changes to database
	if err := database.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update profile"})
		return
	}

	// send success response
	c.JSON(http.StatusOK, ProfileEditResponse{Message: "Profile updated successfully"})
}
