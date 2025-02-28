package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup a test database (In-memory SQLite)
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) // In-memory DB
	if err != nil {
		panic("failed to connect to the database")
	}
	db.AutoMigrate(&models.User{}, &models.UserProfile{}) // Migrate schema
	return db
}

// Test the RegisterUser function
func TestRegisterUser(t *testing.T) {
	db := setupTestDB()
	database.DB = db // Replace the global DB with the test DB

	// Set up Gin for testing
	r := gin.Default()
	r.POST("/auth/register", RegisterUser)

	// Create a registration request body
	requestBody := UserRegistrationRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Perform the POST request
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := performRequest(r, req)

	// Assert the response status and body
	assert.Equal(t, http.StatusCreated, w.Code)

	var response UserRegistrationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Registration successful", response.Message)

	// Check if the user was inserted into the database
	var user models.User
	err = db.Where("email = ?", requestBody.Email).First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Email, requestBody.Email)
}

// Test the LoginUser function
func TestLoginUser(t *testing.T) {
	db := setupTestDB()
	database.DB = db // Replace the global DB with the test DB

	// Set up Gin for testing
	r := gin.Default()
	r.POST("/auth/login", LoginUser)

	// Prepopulate the database with a user
	password, _ := utils.HashPassword("password123")
	user := models.User{
		Email:    "test@example.com",
		Password: password,
	}
	db.Create(&user)

	// Create a login request body
	requestBody := UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// Perform the POST request
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := performRequest(r, req)

	// Assert the response status and body
	assert.Equal(t, http.StatusOK, w.Code)

	var response UserLoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", response.Message)
	assert.NotEmpty(t, response.Token)

	// Check if the user ID is correct in the response
	var dbUser models.User
	err = db.Where("email = ?", requestBody.Email).First(&dbUser).Error
	assert.NoError(t, err)
	assert.Equal(t, dbUser.ID, response.UserID)
}

// Helper function to perform a request and return the response recorder
func performRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
