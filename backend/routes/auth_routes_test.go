package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock RegisterUser function for testing
func MockRegisterUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Mock LoginUser function for testing
func MockLoginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func TestAuthRoutes_RegisterUser(t *testing.T) {
	// Setup the Gin engine for testing
	r := gin.Default()

	// Register the mock function for testing
	r.POST("/auth/register", MockRegisterUser)

	// Create a mock request to test the route
	reqBody := `{"username": "testuser", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a recorder to capture the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])
}

func TestAuthRoutes_LoginUser(t *testing.T) {
	// Setup the Gin engine for testing
	r := gin.Default()

	// Register the mock function for testing
	r.POST("/auth/login", MockLoginUser)

	// Create a mock request to test the route
	reqBody := `{"username": "testuser", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a recorder to capture the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", response["message"])
}
