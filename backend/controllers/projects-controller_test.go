package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	database.DB.AutoMigrate(&models.User{}, &models.UserProfile{}, &models.Project{}, &models.Collaborator{}, &models.Invitation{})

	// Clean up existing data - Note the order matters due to foreign key constraints
	database.DB.Exec("DELETE FROM invitations")
	database.DB.Exec("DELETE FROM collaborators")
	database.DB.Exec("DELETE FROM projects")
	database.DB.Exec("DELETE FROM user_profiles")
	database.DB.Exec("DELETE FROM users")

	// Reset auto-increment counters
	database.DB.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name IN ('users', 'projects', 'collaborators', 'user_profiles', 'invitations')")
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

func TestProjectCollaboration(t *testing.T) {
	setupProjectsTest(t)

	// Create two test users
	owner := models.User{Email: "owner@example.com", Password: "password"}
	collaborator := models.User{Email: "collaborator@example.com", Password: "password"}
	database.DB.Create(&owner)
	database.DB.Create(&collaborator)

	// Generate tokens for both users
	ownerToken, err := utils.GenerateJWT(owner.ID, owner.Email)
	assert.NoError(t, err)
	collaboratorToken, err := utils.GenerateJWT(collaborator.ID, collaborator.Email)
	assert.NoError(t, err)

	// Setup router with middleware
	router := gin.Default()
	router.POST("/projects", middleware.AuthRequired(), controllers.CreateProject)
	router.POST("/projects/:id/collaborators", middleware.AuthRequired(), controllers.InviteCollaborator)
	router.GET("/projects/invitations", middleware.AuthRequired(), controllers.GetProjectInvitations)
	router.POST("/projects/:id/collaborators/invitations/:invitationId/:action", middleware.AuthRequired(), controllers.RespondToProjectInvitation)

	// Step 1: Owner creates a project
	projectReq := map[string]interface{}{
		"title":           "Test Project",
		"description":     "This is a test project",
		"required_skills": []string{"Go", "Testing"},
		"visibility":      "private",
		"status":          "open",
	}
	projectJson, _ := json.Marshal(projectReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(projectJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ownerToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Extract project ID from response
	var projectResponse struct {
		ID uint `json:"id"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &projectResponse)
	assert.NoError(t, err)

	// Step 2: Owner invites collaborator
	inviteReq := map[string]interface{}{
		"email": collaborator.Email,
		"role":  "programmer",
	}
	inviteJson, _ := json.Marshal(inviteReq)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", fmt.Sprintf("/projects/%d/collaborators", projectResponse.ID), bytes.NewBuffer(inviteJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ownerToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Step 3: Collaborator views all their invitations
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/projects/invitations", nil)
	req.Header.Set("Authorization", "Bearer "+collaboratorToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var invitationsResponse struct {
		Invitations []struct {
			ID           uint      `json:"id"`
			ProjectID    uint      `json:"project_id"`
			ProjectTitle string    `json:"project_title"`
			Role         string    `json:"role"`
			Status       string    `json:"status"`
			CreatedAt    time.Time `json:"created_at"`
		} `json:"invitations"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &invitationsResponse)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(invitationsResponse.Invitations))
	assert.Equal(t, "pending", invitationsResponse.Invitations[0].Status)
	assert.Equal(t, "Test Project", invitationsResponse.Invitations[0].ProjectTitle)

	// Step 4: Collaborator accepts invitation
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", fmt.Sprintf("/projects/%d/collaborators/invitations/%d/accept", projectResponse.ID, invitationsResponse.Invitations[0].ID), nil)
	req.Header.Set("Authorization", "Bearer "+collaboratorToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Step 5: Verify collaborator was added
	var collaboratorCount int64
	database.DB.Model(&models.Collaborator{}).Where("project_id = ? AND user_id = ?", projectResponse.ID, collaborator.ID).Count(&collaboratorCount)
	assert.Equal(t, int64(1), collaboratorCount)

	// Step 6: Verify invitation is no longer pending
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/projects/invitations", nil)
	req.Header.Set("Authorization", "Bearer "+collaboratorToken)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &invitationsResponse)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(invitationsResponse.Invitations))
}

func TestGetProjectInvitationsUnauthorized(t *testing.T) {
	setupProjectsTest(t)

	router := gin.Default()
	router.GET("/projects/invitations", middleware.AuthRequired(), controllers.GetProjectInvitations)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/invitations", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetProjectInvitationsNoInvitations(t *testing.T) {
	setupProjectsTest(t)

	// Create a test user
	user := models.User{Email: "test@example.com", Password: "password"}
	database.DB.Create(&user)

	// Generate token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	assert.NoError(t, err)

	router := gin.Default()
	router.GET("/projects/invitations", middleware.AuthRequired(), controllers.GetProjectInvitations)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/projects/invitations", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Invitations []struct{} `json:"invitations"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(response.Invitations))
}

func TestListUserProjects(t *testing.T) {
	setupProjectsTest(t)

	// Create test users
	owner := models.User{Email: "owner@example.com", Password: "password"}
	collaborator := models.User{Email: "collaborator@example.com", Password: "password"}
	otherUser := models.User{Email: "other@example.com", Password: "password"}
	database.DB.Create(&owner)
	database.DB.Create(&collaborator)
	database.DB.Create(&otherUser)

	// Generate tokens
	ownerToken, err := utils.GenerateJWT(owner.ID, owner.Email)
	assert.NoError(t, err)
	collaboratorToken, err := utils.GenerateJWT(collaborator.ID, collaborator.Email)
	assert.NoError(t, err)

	// Create projects
	projects := []models.Project{
		{
			Title:       "Owner's Project",
			Description: "Description 1",
			OwnerID:     owner.ID,
			Visibility:  "private",
			Status:      "open",
		},
		{
			Title:       "Collaborator's Project",
			Description: "Description 2",
			OwnerID:     otherUser.ID,
			Visibility:  "private",
			Status:      "in-progress",
		},
	}

	for i := range projects {
		projects[i].SetRequiredSkills([]string{"Go", "Testing"})
		database.DB.Create(&projects[i])
	}

	// Add collaborator to the second project
	collab := models.Collaborator{
		ProjectID: projects[1].ID,
		UserID:    collaborator.ID,
		Role:      "programmer",
	}
	database.DB.Create(&collab)

	// Setup router
	router := gin.Default()
	router.GET("/projects/user", middleware.AuthRequired(), controllers.ListUserProjects)

	t.Run("Owner sees their own project", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects/user", nil)
		req.Header.Set("Authorization", "Bearer "+ownerToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response controllers.ProjectListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Projects))
		assert.Equal(t, "Owner's Project", response.Projects[0].Title)
	})

	t.Run("Collaborator sees project they're involved in", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects/user", nil)
		req.Header.Set("Authorization", "Bearer "+collaboratorToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response controllers.ProjectListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Projects))
		assert.Equal(t, "Collaborator's Project", response.Projects[0].Title)
	})

	t.Run("Unauthorized access", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/projects/user", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
