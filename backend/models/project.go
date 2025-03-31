package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
)

type Project struct {
	gorm.Model
	Title          string         `gorm:"not null" json:"title"`
	Description    string         `json:"description"`
	OwnerID        uint           `gorm:"not null" json:"owner_id"`
	RequiredSkills string         `gorm:"type:text" json:"required_skills"` // store as JSON string
	Visibility     string         `gorm:"not null" json:"visibility"`
	Status         string         `gorm:"not null" json:"status"`
	Collaborators  []Collaborator `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE;" json:"collaborators"`
}

// Convert RequiredSkills from JSON to []string when reading from DB
func (p *Project) GetRequiredSkills() []string {
	var skills []string
	if err := json.Unmarshal([]byte(p.RequiredSkills), &skills); err != nil {
		log.Println("Error unmarshaling RequiredSkills:", err)
	}
	return skills
}

// Convert []string to JSON before saving RequiredSkills to DB
func (p *Project) SetRequiredSkills(skills []string) {
	skillsJSON, err := json.Marshal(skills)
	if err != nil {
		log.Println("Error marshaling RequiredSkills:", err)
		return
	}
	p.RequiredSkills = string(skillsJSON)
}
