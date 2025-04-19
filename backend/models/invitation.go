package models

import (
	"time"

	"gorm.io/gorm"
)

type InvitationStatus string
type CollaboratorRole string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusRejected InvitationStatus = "rejected"

	CollaboratorRoleProgrammer CollaboratorRole = "programmer"
	CollaboratorRoleEditor     CollaboratorRole = "editor"
	CollaboratorRoleOwner      CollaboratorRole = "owner"
)

type Invitation struct {
	gorm.Model
	ProjectID    uint             `json:"project_id"`
	Project      Project          `json:"-" gorm:"foreignKey:ProjectID"`
	InviterID    uint             `json:"inviter_id"`
	Inviter      User             `json:"-" gorm:"foreignKey:InviterID"`
	Email        string           `json:"email"`
	Role         CollaboratorRole `json:"role"`
	Status       InvitationStatus `json:"status"`
	ResponseDate *time.Time       `json:"response_date,omitempty"`
}
