package dto

import "time"

// ChannelDto represents a channel in API responses
type ChannelDto struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ProjectID string    `json:"project_id"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChannelCreateRequest represents a request to create a new channel
type ChannelCreateRequest struct {
	Name      string `json:"name" binding:"required"`
	ProjectID string `json:"project_id" binding:"required"`
	IsPrivate bool   `json:"is_private"` // If true, creator will be added as member
}

// ChannelUpdateRequest represents a request to update a channel
type ChannelUpdateRequest struct {
	Name *string `json:"name,omitempty"`
}
