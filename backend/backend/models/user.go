package models

import (
	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	// Add other fields as needed (e.g., Address, Phone)
}
