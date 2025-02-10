package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database object
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var err error
	// Replace with your DB credentials
	DB, err = gorm.Open(mysql.Open("user:password@/dbname"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
}
