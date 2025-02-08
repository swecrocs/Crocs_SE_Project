package database

import (
	"backend/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// global database instance
var DB *gorm.DB

// initialize database connection
func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to users database: ", err)
	}

	DB.AutoMigrate(&models.User{}, &models.UserProfile{})
}
