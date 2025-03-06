package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"backend/controllers"
	"backend/database"
	"backend/middleware"
	"backend/models"
)

func setupUsersTest(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the test database: %v", err)
	}

	// Run migrations or setup test data here if needed
	database.DB.AutoMigrate(&models.User{}, &models.UserProfile{})
}

func registerAndLoginUser(t *testing.T, email string) (string, uint) {
	setupUsersTest(t) // Initialize the database for this test

	// Register a new user
	router := gin.Default()
	router.POST("/auth/register", controllers.RegisterUser)

	w := httptest.NewRecorder()
	reqBody := `{"email": "` + email + `", "password": "securepassword"}`
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Login the user to get a JWT token
	router.POST("/auth/login", controllers.LoginUser)

	w = httptest.NewRecorder()
	reqBody = `{"email": "` + email + `", "password": "securepassword"}`
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Extract the token from the response
	var loginResponse struct {
		Token  string `json:"token"`
		UserID uint   `json:"user_id"`
	}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &loginResponse))

	return loginResponse.Token, loginResponse.UserID
}

func TestRetrieveUserProfile(t *testing.T) {
	setupUsersTest(t) // Initialize the database for this test

	// Create a user directly in the database
	user := models.User{Email: "testuser2@example.com", Password: "securepassword"}
	database.DB.Create(&user)

	// Create a user profile directly in the database
	userProfile := models.UserProfile{UserID: user.ID, FullName: "Test User", Bio: "This is a test bio", Affiliation: "Test Organization"}
	database.DB.Create(&userProfile)

	router := gin.Default()
	router.GET("/users/:id/profile", controllers.RetrieveUserProfile)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/profile", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user_id")
}

func TestEditUserProfile(t *testing.T) {
	token, _ := registerAndLoginUser(t, "testuser1@example.com")

	router := gin.Default()
	router.PUT("/users/:id/profile", controllers.EditUserProfile)

	w := httptest.NewRecorder()
	reqBody := `{"full_name": "Test User", "bio": "This is a test bio", "affiliation": "Test Organization"}`
	req, _ := http.NewRequest("PUT", "/users/1/profile", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Profile updated successfully")
}

func TestEditUserProfileUnauthorized(t *testing.T) {
	setupUsersTest(t) // Initialize the database for this test

	// Register and login the first user
	token1, _ := registerAndLoginUser(t, "testuser3@example.com")

	// Register and login a second user
	_, userID2 := registerAndLoginUser(t, "testuser4@example.com")

	// Attempt to edit the second user's profile with the first user's token
	router := gin.Default()
	router.PUT("/users/:id/profile", middleware.AuthRequired(), middleware.SameUserOnly(), controllers.EditUserProfile)

	w := httptest.NewRecorder()
	reqBody := `{"full_name": "Unauthorized User", "bio": "This should not be allowed", "affiliation": "Unauthorized Org"}`
	req, _ := http.NewRequest("PUT", "/users/"+strconv.Itoa(int(userID2))+"/profile", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1) // Use the first user's token

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "You can only modify your own profile")
}
