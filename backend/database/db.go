package database

import (
	"backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	// Auto migrate the models
	DB.AutoMigrate(&models.User{}, &models.UserProfile{})
}
