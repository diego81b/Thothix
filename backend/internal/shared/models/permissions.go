package models

import (
	"gorm.io/gorm"

	commonModels "thothix-backend/internal/common/models"
)

// RoleType defines the available role types
type RoleType string

const (
	// Simplified role system
	RoleAdmin    RoleType = "admin"    // Can manage everything
	RoleManager  RoleType = "manager"  // Can manage everything except users
	RoleUser     RoleType = "user"     // Can participate in assigned projects/channels, create 1:1 chats
	RoleExternal RoleType = "external" // Can only participate in public channels
)

// Permission defines specific permissions
type Permission string

const (
	// User management permissions
	PermissionUserManage Permission = "user:manage"

	// Project permissions
	PermissionProjectCreate Permission = "project:create"
	PermissionProjectRead   Permission = "project:read"
	PermissionProjectUpdate Permission = "project:update"
	PermissionProjectDelete Permission = "project:delete"
	PermissionProjectManage Permission = "project:manage"

	// Channel permissions
	PermissionChannelCreate       Permission = "channel:create"
	PermissionChannelRead         Permission = "channel:read"
	PermissionChannelUpdate       Permission = "channel:update"
	PermissionChannelDelete       Permission = "channel:delete"
	PermissionChannelManage       Permission = "channel:manage"
	PermissionChannelReadAssigned Permission = "channel:read_assigned" // Read assigned channels

	// Message permissions
	PermissionMessageCreate Permission = "message:create"
	PermissionMessageRead   Permission = "message:read"
	PermissionMessageUpdate Permission = "message:update"
	PermissionMessageDelete Permission = "message:delete"

	// Direct message permissions
	PermissionDMCreate Permission = "dm:create" // Create 1:1 conversations

	// File permissions
	PermissionFileUpload Permission = "file:upload"
	PermissionFileRead   Permission = "file:read"
	PermissionFileDelete Permission = "file:delete"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[RoleType][]Permission{
	// Admin: can manage everything
	RoleAdmin: {
		PermissionUserManage,
		PermissionProjectCreate, PermissionProjectRead, PermissionProjectUpdate, PermissionProjectDelete, PermissionProjectManage,
		PermissionChannelCreate, PermissionChannelRead, PermissionChannelUpdate, PermissionChannelDelete, PermissionChannelManage, PermissionChannelReadAssigned,
		PermissionMessageCreate, PermissionMessageRead, PermissionMessageUpdate, PermissionMessageDelete,
		PermissionDMCreate,
		PermissionFileUpload, PermissionFileRead, PermissionFileDelete,
	},

	// Manager: can manage everything except users
	RoleManager: {
		PermissionProjectCreate, PermissionProjectRead, PermissionProjectUpdate, PermissionProjectDelete, PermissionProjectManage,
		PermissionChannelCreate, PermissionChannelRead, PermissionChannelUpdate, PermissionChannelDelete, PermissionChannelManage, PermissionChannelReadAssigned,
		PermissionMessageCreate, PermissionMessageRead, PermissionMessageUpdate, PermissionMessageDelete,
		PermissionDMCreate,
		PermissionFileUpload, PermissionFileRead, PermissionFileDelete,
	},

	// User: can participate in assigned projects/channels, create 1:1 chats
	RoleUser: {
		PermissionProjectRead, // Can only read projects they're assigned to
		PermissionChannelRead, PermissionChannelReadAssigned,
		PermissionMessageCreate, PermissionMessageRead, PermissionMessageUpdate,
		PermissionDMCreate,
		PermissionFileUpload, PermissionFileRead,
	},

	// External: can only participate in public channels
	RoleExternal: {
		PermissionChannelRead,                          // Only public channels (controlled by hasChannelAccess)
		PermissionMessageCreate, PermissionMessageRead, // Only in public channels
		PermissionFileUpload, PermissionFileRead,
	},
}

// HasPermission checks if a role has a specific permission
func (r RoleType) HasPermission(permission Permission) bool {
	permissions, exists := RolePermissions[r]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// This will need to be refactored to use domain models after we complete the migration
// For now, we'll keep these as interfaces that can be implemented by domain models

// UserProvider interface for getting user data
type UserProvider interface {
	GetSystemRole() RoleType
}

// ChannelProvider interface for getting channel data
type ChannelProvider interface {
	GetIsPrivate() bool
}

// GetUserRole gets the user's system role from database
// This is a temporary implementation that uses the old models
// Will be updated once all models are migrated to domains
func GetUserRole(db *gorm.DB, userID string) (RoleType, error) {
	// This will be updated to use usersDomain.User once migration is complete
	var result struct {
		SystemRole RoleType `json:"system_role"`
	}

	if err := db.Table("users").Select("system_role").Where("id = ?", userID).First(&result).Error; err != nil {
		return RoleUser, err // Default to user role on error
	}
	return result.SystemRole, nil
}

// HasUserPermission checks if a user has a specific permission
func HasUserPermission(db *gorm.DB, userID string, permission Permission, resourceType, resourceID *string) bool {
	// Get user's system role
	userRole, err := GetUserRole(db, userID)
	if err != nil {
		return false
	}

	// Check if the role has the permission
	if userRole.HasPermission(permission) {
		// For channel-specific permissions, check additional constraints
		if resourceType != nil && *resourceType == "channel" && resourceID != nil {
			return hasChannelAccess(db, userID, *resourceID, permission)
		}
		// For project-specific permissions, check additional constraints
		if resourceType != nil && *resourceType == "project" && resourceID != nil {
			return hasProjectAccess(db, userID, *resourceID, permission)
		}
		return true
	}

	return false
}

// hasChannelAccess checks if user has access to a specific channel
// Temporary implementation using raw SQL queries
func hasChannelAccess(db *gorm.DB, userID, channelID string, permission Permission) bool {
	// Get channel info
	var channel struct {
		ProjectID string `json:"project_id"`
	}
	if err := db.Table("channels").Select("project_id").Where("id = ?", channelID).First(&channel).Error; err != nil {
		return false
	}

	// Check if channel is private (has project_id)
	isPrivate := channel.ProjectID != ""

	// Get user's system role
	userRole, err := GetUserRole(db, userID)
	if err != nil {
		return false
	}

	// External users can only access public channels
	if userRole == RoleExternal {
		return !isPrivate
	}

	// For private channels, check if user is a member or has elevated role
	if isPrivate {
		if userRole == RoleAdmin || userRole == RoleManager {
			return true
		}
		// Check if user is a member of the channel
		var count int64
		db.Table("channel_members").Where("channel_id = ? AND user_id = ?", channelID, userID).Count(&count)
		return count > 0
	}

	// Public channels are accessible to all authenticated users
	return true
}

// hasProjectAccess checks if user has access to a specific project
func hasProjectAccess(db *gorm.DB, userID, projectID string, permission Permission) bool {
	userRole, err := GetUserRole(db, userID)
	if err != nil {
		return false
	}

	// Admins and managers have access to all projects
	if userRole == RoleAdmin || userRole == RoleManager {
		return true
	}

	// For regular users and external users, check if they are project members
	var count int64
	db.Table("project_members").Where("project_id = ? AND user_id = ?", projectID, userID).Count(&count)
	return count > 0
}

// Legacy function for backward compatibility
func HasUserPermissionSimple(userID string, permission Permission) bool {
	// This function is deprecated - use HasUserPermission with db parameter instead
	return false
}

// UserRole represents a user's role assignment
// swagger:model UserRole
type UserRole struct {
	commonModels.BaseModel
	ResourceID   *string  `json:"resource_id,omitempty"`   // For project/channel specific roles (not used in simplified system)
	ResourceType *string  `json:"resource_type,omitempty"` // "project", "channel", null for system roles (not used in simplified system)
	UserID       string   `json:"user_id"`
	Role         RoleType `json:"role"`
}

// TableName specifies the table name for the UserRole model
func (UserRole) TableName() string {
	return "user_roles"
}
