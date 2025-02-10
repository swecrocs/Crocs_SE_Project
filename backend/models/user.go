package models

import (
	"backend/database" // Make sure you're importing the correct package
	"fmt"
)

// User struct representing a user in the database
type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

// GetUserByEmail retrieves a user by email from the database
func GetUserByEmail(email string) (User, error) {
	var user User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, fmt.Errorf("user not found")
	}
	return user, nil
}
