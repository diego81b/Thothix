package domain

import (
	commonModels "thothix-backend/internal/common/models"
)

// Channel represents a chat channel in the chat domain
type Channel struct {
	commonModels.BaseModel
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
	IsPrivate bool   `json:"is_private" gorm:"-"` // Computed field, not stored in DB
}

// LoadIsPrivate calculates and sets the IsPrivate field based on channel members
func (c *Channel) LoadIsPrivate(db interface{}) error {
	// Simplified implementation - assume private if has project_id
	c.IsPrivate = c.ProjectID != ""
	return nil
}

// ChannelMember represents a user's membership in a channel
type ChannelMember struct {
	commonModels.BaseModel
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
}
