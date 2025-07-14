package dto

import "time"

// ProjectDto represents a project in API responses
type ProjectDto struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectCreateRequest represents a request to create a new project
type ProjectCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// ProjectUpdateRequest represents a request to update a project
type ProjectUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
