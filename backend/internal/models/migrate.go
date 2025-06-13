package models

import "gorm.io/gorm"

// AutoMigrate all tables and foreign keys
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Project{},
		&ProjectMember{},
		&Channel{},
		&ChannelMember{},
		&Message{},
		&File{},
		&UserRole{},
	)
}
