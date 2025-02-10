package controllers

import (
	"backend/models"
	"backend/utils"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// LoginUser godoc
// @Summary      Authenticate user
// @Description  Authenticate user credentials and return a JWT token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body models.User true "User login details"
// @Success      200 {object} map[string]string
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} utils.ErrorResponse
// @Router       /auth/login [post]
func LoginUser(c *gin.Context) {
	var loginRequest models.User

	// Bind JSON request
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "Invalid request"})
		return
	}

	// Fetch user from database
	user, err := models.GetUserByEmail(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	// Return token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
