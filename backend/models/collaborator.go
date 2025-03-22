package models

import "gorm.io/gorm"

type Collaborator struct {
	gorm.Model
	ProjectID uint   `gorm:"not null;index" json:"project_id"`
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	Role      string `gorm:"not null" json:"role"`
}
