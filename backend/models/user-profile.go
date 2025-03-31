package models

import "gorm.io/gorm"

type UserProfile struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"unique"`
	FullName    string `json:"full_name"`
	Bio         string `json:"bio"`
	Affiliation string `json:"affiliation"`
	Skills      string `json:"skills"`
	Role        string `json:"role"`
	Projects    string `json:"projects"`
}
