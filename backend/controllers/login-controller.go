package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// GenerateJWT generates a JWT token for the user after successful login
func GenerateJWT(user models.User) (string, error) {
	// Define token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create JWT claims, which includes the username and expiry time
	claims := &jwt.Claims{
		Subject:   string(user.ID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	// Create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with a secret key (make sure to store the secret safely)
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// LoginUser godoc
// @Summary      Log in a user
// @Description  Authenticate user by email and password and generate a JWT token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        requestBody body UserLoginRequest true "User credentials"
// @Success      200 {object} UserLoginResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /auth/login [post]
func LoginUser(c *gin.Context) {
	var loginRequest UserLoginRequest

	// Bind JSON input
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	// Retrieve user by email
	var user models.User
	if err := database.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Check if the password is correct
	if err := utils.CheckPasswordHash(loginRequest.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not generate token"})
		return
	}

	// Send the response with token
	c.JSON(http.StatusOK, UserLoginResponse{
		Message: "Login successful",
		Token:   token,
	})
}
