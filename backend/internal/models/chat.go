package models

import "gorm.io/gorm"

// Channel represents a chat channel
// swagger:model Channel
type Channel struct {
	BaseModel
	Name      string  `json:"name"`
	ProjectID string  `json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"project,omitempty"`
	Members   []User  `gorm:"many2many:channel_members" json:"members,omitempty"`
	IsPrivate bool    `json:"is_private" gorm:"-"` // Computed field, not stored in DB
}

// LoadIsPrivate calculates and sets the IsPrivate field based on channel members
func (c *Channel) LoadIsPrivate(db *gorm.DB) error {
	var count int64
	err := db.Model(&ChannelMember{}).Where("channel_id = ?", c.ID).Count(&count).Error
	if err != nil {
		return err
	}
	c.IsPrivate = count > 0
	return nil
}

// ChannelMember represents membership of a user in a channel
// swagger:model ChannelMember
type ChannelMember struct {
	BaseModel
	ChannelID string  `json:"channel_id"`
	Channel   Channel `gorm:"foreignKey:ChannelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"channel,omitempty"`
	UserID    string  `json:"user_id"`
	User      User    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}
