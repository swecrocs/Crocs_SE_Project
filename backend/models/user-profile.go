package models

import (
	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	UserID uint   `gorm:"uniqueIndex"`
	Name   string `gorm:"size:255"`
	Age    int
}
