package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	// Clean up existing data - Note the order matters due to foreign key constraints
	database.DB.Exec("DELETE FROM collaborators")
	database.DB.Exec("DELETE FROM projects")
	database.DB.Exec("DELETE FROM user_profiles")
	database.DB.Exec("DELETE FROM users")

	// Reset auto-increment counters
	database.DB.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name IN ('users', 'projects', 'collaborators', 'user_profiles')")
}

func TestRetrieveProject(t *testing.T) {
	setupProjectsTest(t)

	// Create a test user
	user := models.User{Email: "retrieve_project@example.com", Password: "password"}
	result := database.DB.Create(&user)
	assert.NoError(t, result.Error)

	// Create a test project
	project := models.Project{
		Title:       "Test Project",
		Description: "Test Description",
		OwnerID:     user.ID,
		Visibility:  "private",
		Status:      "open",
	}
	project.SetRequiredSkills([]string{"Go", "Testing"})
	result = database.DB.Create(&project)
	assert.NoError(t, result.Error)

	// Setup router
	router := gin.Default()
	router.GET("/projects/:id", controllers.RetrieveProject)

	t.Run("Successful project retrieval", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d", project.ID), nil)
		router.ServeHTTP(w, req)

		// Assert response code
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var response controllers.ProjectRetrievalResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify project details
		assert.Equal(t, project.ID, response.ID)
		assert.Equal(t, project.Title, response.Title)
		assert.Equal(t, project.Description, response.Description)
	})

	t.Run("Project not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects/999", nil)
		router.ServeHTTP(w, req)

		// Assert response code
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Parse error response
		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Error, "Project not found")
	})

	t.Run("Invalid project ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects/invalid", nil)
		router.ServeHTTP(w, req)

		// Assert response code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Parse error response
		var response controllers.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Error, "Failed to fetch project")
	})
}

func TestListProjects(t *testing.T) {
    setupProjectsTest(t)

    // Create some test users
    user1 := models.User{Email: "user1@example.com", Password: "password"}
    user2 := models.User{Email: "user2@example.com", Password: "password"}
    database.DB.Create(&user1)
    database.DB.Create(&user2)

    // Create test projects
    projects := []models.Project{
        {
            Title:       "Project 1",
            Description: "Description 1",
            OwnerID:     user1.ID,
            Visibility:  "private",
            Status:      "open",
        },
        {
            Title:       "Project 2",
            Description: "Description 2",
            OwnerID:     user2.ID,
            Visibility:  "private",
            Status:      "in-progress",
        },
    }

    for i := range projects {
        projects[i].SetRequiredSkills([]string{"Go", "Testing"})
        database.DB.Create(&projects[i])
    }

    // Setup router
    router := gin.Default()
    router.GET("/projects", controllers.ListProjects)

    // Make request
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/projects", nil)
    router.ServeHTTP(w, req)

    // Verify response
    assert.Equal(t, http.StatusOK, w.Code)

    var response controllers.ProjectListResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)

    // Verify projects count
    assert.Equal(t, 2, len(response.Projects))

    // Verify project details
    assert.Equal(t, projects[0].Title, response.Projects[0].Title)
    assert.Equal(t, projects[0].Description, response.Projects[0].Description)
    assert.Equal(t, projects[0].OwnerID, response.Projects[0].OwnerID)
    assert.Equal(t, projects[1].Title, response.Projects[1].Title)
    assert.Equal(t, projects[1].Description, response.Projects[1].Description)
    assert.Equal(t, projects[1].OwnerID, response.Projects[1].OwnerID)
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
