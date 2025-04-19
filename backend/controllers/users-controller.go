package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProfileRetrievalResponse struct {
	UserID      uint   `json:"user_id"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Bio         string `json:"bio"`
	Affiliation string `json:"affiliation"`
	Skills      string `json:"skills"`
	Role        string `json:"role"`
	Projects    string `json:"projects"`
}

// RetrieveUserProfile godoc
// @Summary      Get user profile
// @Description  Retrieve user profile information by user ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  ProfileRetrievalResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /users/{id}/profile [get]
func RetrieveUserProfile(c *gin.Context) {
	userID := c.Param("id")

	// get user and user profile in a single query
	var user models.User
	if err := database.DB.Preload("Profile").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	// create response with new fields
	response := ProfileRetrievalResponse{
		UserID:      user.ID,
		Email:       user.Email,
		FullName:    user.Profile.FullName,
		Bio:         user.Profile.Bio,
		Affiliation: user.Profile.Affiliation,
		Skills:      user.Profile.Skills,
		Role:        user.Profile.Role,
		Projects:    user.Profile.Projects,
	}

	// respond on success
	c.JSON(http.StatusOK, response)
}

type ProfileEditRequest struct {
	FullName    string `json:"full_name"`
	Bio         string `json:"bio"`
	Affiliation string `json:"affiliation"`
	Skills      string `json:"skills"`
	Role        string `json:"role"`
	Projects    string `json:"projects"`
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
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Param        request body ProfileEditRequest true "Profile information"
// @Success      200 {object} ProfileEditResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse "Unauthorized - Missing or invalid JWT token"
// @Failure      403 {object} ErrorResponse "Forbidden - Cannot modify another user's profile"
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /users/{id}/profile [put]
func EditUserProfile(c *gin.Context) {
	userID := c.Param("id")

	// convert userID string to uint
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// bind request JSON to ProfileEditRequest struct
	var request ProfileEditRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	// find or create the profile linked to this user
	var profile models.UserProfile
	result := database.DB.Where("user_id = ?", userID).FirstOrCreate(&profile, models.UserProfile{UserID: uint(uid)})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to find/create profile"})
		return
	}

	// update profile fields with new fields
	profile.FullName = request.FullName
	profile.Bio = request.Bio
	profile.Affiliation = request.Affiliation
	profile.Skills = request.Skills
	profile.Role = request.Role
	profile.Projects = request.Projects

	// save changes to database
	if err := database.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update profile"})
		return
	}

	// send success response
	c.JSON(http.StatusOK, ProfileEditResponse{Message: "Profile updated successfully"})
}
