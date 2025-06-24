package models

import (
	"log"

	"gorm.io/gorm"
)

// AutoMigrate all tables and foreign keys
func AutoMigrate(db *gorm.DB) error {
	log.Println("🔄 Starting database migration...")

	err := db.AutoMigrate(
		&User{},
		&Project{},
		&ProjectMember{},
		&Channel{},
		&ChannelMember{},
		&Message{},
		&File{},
		&UserRole{},
	)
	if err != nil {
		log.Printf("❌ Migration failed: %v", err)
		return err
	}

	log.Println("✅ Database migration completed successfully")
	return nil
}
