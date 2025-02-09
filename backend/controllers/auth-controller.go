package controllers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key") // Change this key in production

// LoginRequest model for login data
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse model to return JWT token
type LoginResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary      User login
// @Description  Login with email and password to get a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        login body LoginRequest true "User login credentials"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	// Compare the password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// Helper function to generate JWT token
func generateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
