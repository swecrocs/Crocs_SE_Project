package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"backend/controllers"
	"backend/database"
	"backend/middleware"
	"backend/models"
	"backend/utils"
)

func setupProjectsTest(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the test database: %v", err)
	}

	// Run migrations
	database.DB.AutoMigrate(&models.User{}, &models.UserProfile{}, &models.Project{}, &models.Collaborator{})
}

func TestCreateProject(t *testing.T) {
	setupProjectsTest(t)

	// Create a test user directly
	user := models.User{Email: "project_user@example.com", Password: "password"}
	database.DB.Create(&user)

	// Generate a token for the user
	token, err := utils.GenerateJWT(user.ID, user.Email)
	assert.NoError(t, err)

	// Setup router with middleware
	router := gin.Default()
	router.POST("/projects", middleware.AuthRequired(), controllers.CreateProject)

	// Create project request
	reqBody := map[string]interface{}{
		"title":           "Test Project",
		"description":     "This is a test project",
		"required_skills": []string{"Go", "Testing"},
		"visibility":      "private",
		"status":          "open",
	}
	jsonData, _ := json.Marshal(reqBody)

	// Create the request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Set the auth token

	// Serve the request
	router.ServeHTTP(w, req)

	// Verify the response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Project successfully created")

	// Verify that the project was created in the database
	var count int64
	database.DB.Model(&models.Project{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// Verify that the collaborator was added
	var collaborator models.Collaborator
	err = database.DB.First(&collaborator).Error
	assert.NoError(t, err)
	assert.Equal(t, user.ID, collaborator.UserID)
	assert.Equal(t, "owner", collaborator.Role)
}

func TestCreateProjectUnauthorized(t *testing.T) {
	setupProjectsTest(t)

	// Setup router with middleware
	router := gin.Default()
	router.POST("/projects", middleware.AuthRequired(), controllers.CreateProject)

	// Create project request
	reqBody := map[string]interface{}{
		"title":           "Test Project",
		"description":     "This is a test project",
		"required_skills": []string{"Go", "Testing"},
		"visibility":      "private",
		"status":          "open",
	}
	jsonData, _ := json.Marshal(reqBody)

	// Create the request without auth token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Serve the request
	router.ServeHTTP(w, req)

	// Verify the response shows unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}