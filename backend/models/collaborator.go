package models

import "gorm.io/gorm"

type Collaborator struct {
	gorm.Model
	ProjectID uint        `gorm:"not null;index" json:"project_id"`
	UserID    uint        `gorm:"not null;index" json:"user_id"`
	User      UserProfile `gorm:"foreignKey:UserID" json:"user"`
	Role      string      `gorm:"not null" json:"role"`
	Status    string      `json:"status"`
}
