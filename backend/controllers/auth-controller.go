package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegistrationResponse struct {
	Message string `json:"message" example:"Registration successful"`
	UserID  uint   `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Create a new user account using user credentials. The provided password is hashed before storing to database. A blank user profile is created.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        requestBody body UserRegistrationRequest true "User credentials"
// @Success      201 {object} UserRegistrationResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/register [post]
func RegisterUser(c *gin.Context) {
	var requestBody UserRegistrationRequest

	// validate input request against expected schema
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	// validate email formatting
	if !utils.IsEmailValid(requestBody.Email) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid email format"})
	}

	// hash password
	hashedPassword, err := utils.HashPassword(requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
		return
	}

	// create new user
	user := models.User{Email: requestBody.Email, Password: hashedPassword}
	// ensure email does not already exist in user database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Email already registered"})
		return
	}

	// create blank user profile
	userProfile := models.UserProfile{UserID: user.ID}
	if err := database.DB.Create(&userProfile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user profile"})
		return
	}

	// respond upon successful registration
	c.JSON(
		http.StatusCreated,
		UserRegistrationResponse{Message: "Registration successful", UserID: user.ID},
	)
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
	Token   string `json:"token"`
}

// LoginUser godoc
// @Summary      Login user
// @Description  Authenticate user login with email and password, returns JWT token on success.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        requestBody body UserLoginRequest true "User credentials"
// @Success      200 {object} UserLoginResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/login [post]
func LoginUser(c *gin.Context) {
	var requestBody UserLoginRequest

	// validate request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	// find user by email
	var user models.User
	if err := database.DB.Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email"})
		return
	}

	// verify password
	if !utils.CheckPassword(user.Password, requestBody.Password) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid password"})
		return
	}

	// generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
		return
	}

	// return success response with token
	c.JSON(http.StatusOK, UserLoginResponse{
		Message: "Login successful",
		UserID:  user.ID,
		Token:   token,
	})
}
