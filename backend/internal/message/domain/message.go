package domain

import (
	commonModels "thothix-backend/internal/common/models"
)

// Message represents a chat or direct message in the message domain
type Message struct {
	commonModels.BaseModel
	SenderID   string  `json:"sender_id"`
	ChannelID  *string `json:"channel_id,omitempty"`
	ReceiverID *string `json:"receiver_id,omitempty"`
	Content    string  `json:"content"`
}

// File represents a file uploaded in a message or project
type File struct {
	commonModels.BaseModel
	MessageID *string `json:"message_id,omitempty"`
	ProjectID *string `json:"project_id,omitempty"` // Optional, can be null for files in direct messages
	URL       string  `json:"url"`
}
