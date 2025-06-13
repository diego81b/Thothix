package models

import "gorm.io/gorm"

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
		PermissionChannelCreate, PermissionChannelRead, PermissionChannelUpdate, PermissionChannelDelete, PermissionChannelManage,
		PermissionChannelReadAssigned,
		PermissionMessageCreate, PermissionMessageRead, PermissionMessageUpdate, PermissionMessageDelete,
		PermissionDMCreate,
		PermissionFileUpload, PermissionFileRead, PermissionFileDelete,
	},

	// Manager: can manage everything except users
	RoleManager: {
		PermissionProjectCreate, PermissionProjectRead, PermissionProjectUpdate, PermissionProjectDelete, PermissionProjectManage,
		PermissionChannelCreate, PermissionChannelRead, PermissionChannelUpdate, PermissionChannelDelete, PermissionChannelManage,
		PermissionChannelReadAssigned,
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
		PermissionFileRead, // Only read files
	},
}

// UserRole represents a user's role assignment
// swagger:model UserRole
type UserRole struct {
	BaseModel
	UserID       string   `json:"user_id"`
	User         User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Role         RoleType `json:"role"`
	ResourceID   *string  `json:"resource_id,omitempty"`   // For project/channel specific roles (not used in simplified system)
	ResourceType *string  `json:"resource_type,omitempty"` // "project", "channel", null for system roles (not used in simplified system)
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

// GetUserRole gets the user's system role from database
func GetUserRole(db *gorm.DB, userID string) (RoleType, error) {
	var user User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return RoleUser, err // Default to user role on error
	}
	return user.SystemRole, nil
}

// HasUserPermission checks if a user has a specific permission
func HasUserPermission(db *gorm.DB, userID string, permission Permission, resourceType *string, resourceID *string) bool {
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
func hasChannelAccess(db *gorm.DB, userID string, channelID string, permission Permission) bool {
	var channel Channel
	if err := db.Where("id = ?", channelID).First(&channel).Error; err != nil {
		return false
	}

	// Load the IsPrivate field
	if err := channel.LoadIsPrivate(db); err != nil {
		return false
	}

	// Get user's system role
	userRole, err := GetUserRole(db, userID)
	if err != nil {
		return false
	}

	// External users can only access public channels
	if userRole == RoleExternal {
		return !channel.IsPrivate
	}

	// For private channels, check if user is a member or has elevated role
	if channel.IsPrivate {
		if userRole == RoleAdmin || userRole == RoleManager {
			return true
		}
		// Check if user is a member of the channel
		var member ChannelMember
		err := db.Where("channel_id = ? AND user_id = ?", channelID, userID).First(&member).Error
		return err == nil
	}

	// Public channels are accessible to all authenticated users
	return true
}

// hasProjectAccess checks if user has access to a specific project
func hasProjectAccess(db *gorm.DB, userID string, projectID string, permission Permission) bool {
	userRole, err := GetUserRole(db, userID)
	if err != nil {
		return false
	}

	// Admins and managers have access to all projects
	if userRole == RoleAdmin || userRole == RoleManager {
		return true
	}

	// For regular users and external users, check if they are project members
	var member ProjectMember
	err = db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&member).Error
	return err == nil
}

// Legacy function for backward compatibility
func HasUserPermissionSimple(userID string, permission Permission) bool {
	// This function is deprecated - use HasUserPermission with db parameter instead
	return false
}
