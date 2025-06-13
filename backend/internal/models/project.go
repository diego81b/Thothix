package models

import (
	"time"
)

// Project represents a project entity
// swagger:model Project
type Project struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProjectMember represents membership of a user in a project
// swagger:model ProjectMember
type ProjectMember struct {
	BaseModel
	UserID    string    `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	ProjectID string    `json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"project,omitempty"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`
}
