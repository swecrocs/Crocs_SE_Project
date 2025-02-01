package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"backend/database"
	"backend/models"
	"backend/utils"
)

func RegisterUser(c *gin.Context) {
	var requestBody struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// create new user
	user := models.User{Email: requestBody.Email, Password: hashedPassword}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}