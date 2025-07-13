package dto

import "time"

// MessageDto represents a message in API responses
type MessageDto struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ChannelID  *string   `json:"channel_id,omitempty"`
	ReceiverID *string   `json:"receiver_id,omitempty"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MessageCreateRequest represents a request to create a new message
type MessageCreateRequest struct {
	ChannelID  *string `json:"channel_id,omitempty"`
	ReceiverID *string `json:"receiver_id,omitempty"`
	Content    string  `json:"content" binding:"required"`
}

// MessageUpdateRequest represents a request to update a message
type MessageUpdateRequest struct {
	Content *string `json:"content,omitempty"`
}

// DirectMessageCreateRequest represents a request to create a direct message
type DirectMessageCreateRequest struct {
	Content     string `json:"content" binding:"required"`
	RecipientID string `json:"recipient_id" binding:"required"`
}
