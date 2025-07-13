package domain

import (
	"time"

	commonModels "thothix-backend/internal/common/models"
)

// Project represents a project entity in the project domain
type Project struct {
	commonModels.BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProjectMember represents membership of a user in a project
type ProjectMember struct {
	commonModels.BaseModel
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`
	UserID    string    `json:"user_id"`
	ProjectID string    `json:"project_id"`
	Role      string    `json:"role"`
}
