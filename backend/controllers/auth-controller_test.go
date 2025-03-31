package controllers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"backend/controllers"
	"backend/database"
	"backend/models"
)

func setupAuthTest(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the test database: %v", err)
	}

	// Run migrations or setup test data here if needed
	database.DB.AutoMigrate(&models.User{}, &models.UserProfile{})
}

func TestRegisterUser(t *testing.T) {
	setupAuthTest(t) // Initialize the database for this test

	router := gin.Default()
	router.POST("/auth/register", controllers.RegisterUser)

	w := httptest.NewRecorder()
	reqBody := `{"email": "testuser@example.com", "password": "securepassword"}`
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Registration successful")
}

func TestLoginUser(t *testing.T) {
	setupAuthTest(t) // Initialize the database for this test

	router := gin.Default()
	router.POST("/auth/login", controllers.LoginUser)

	w := httptest.NewRecorder()
	reqBody := `{"email": "testuser@example.com", "password": "securepassword"}`
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
}
