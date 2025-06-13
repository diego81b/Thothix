package models

// Message represents a chat or direct message
// swagger:model Message
type Message struct {
	BaseModel
	SenderID   string   `json:"sender_id"`
	Sender     User     `gorm:"foreignKey:SenderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"sender,omitempty"`
	ChannelID  *string  `json:"channel_id,omitempty"`
	Channel    *Channel `gorm:"foreignKey:ChannelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"channel,omitempty"`
	ReceiverID *string  `json:"receiver_id,omitempty"`
	Receiver   *User    `gorm:"foreignKey:ReceiverID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"receiver,omitempty"`
	Content    string   `json:"content"`
}

// File represents a file uploaded in a message or project
// swagger:model File
type File struct {
	BaseModel
	URL       string   `json:"url"`
	MessageID *string  `json:"message_id,omitempty"`
	Message   *Message `gorm:"foreignKey:MessageID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"message,omitempty"`
	ProjectID *string  `json:"project_id,omitempty"` // Optional, can be null for files in direct messages
	Project   *Project `gorm:"foreignKey:ProjectID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"project,omitempty"`
}
