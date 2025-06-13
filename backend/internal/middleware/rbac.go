package middleware

import (
	"net/http"
	"thothix-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequirePermission middleware to check if user has specific permission
func RequirePermission(db *gorm.DB, permission models.Permission, resourceType *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("clerk_user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get resource ID from URL params if needed
		var resourceID *string
		if resourceType != nil {
			if id := c.Param("id"); id != "" {
				resourceID = &id
			}
		}

		// Check permission
		if !models.HasUserPermission(db, userID.(string), permission, resourceType, resourceID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSystemRole middleware to check if user has specific system role
func RequireSystemRole(db *gorm.DB, role models.RoleType) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("clerk_user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get user's system role from database
		userRole, err := models.GetUserRole(db, userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user role"})
			c.Abort()
			return
		}

		// Check if user has the required role or higher
		if !hasRequiredSystemRole(userRole, role) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role privileges"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireProjectAccess middleware to check if user can access a project
func RequireProjectAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("clerk_user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		projectID := c.Param("id")
		if projectID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID required"})
			c.Abort()
			return
		}

		// Check if user has access to the project
		resourceType := "project"
		if !models.HasUserPermission(db, userID.(string), models.PermissionProjectRead, &resourceType, &projectID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to project"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireChannelAccess middleware to check if user can access a channel
func RequireChannelAccess(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("clerk_user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		channelID := c.Param("id")
		if channelID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Channel ID required"})
			c.Abort()
			return
		}

		// Check if user has access to the channel
		resourceType := "channel"
		if !models.HasUserPermission(db, userID.(string), models.PermissionChannelRead, &resourceType, &channelID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to channel"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRequiredSystemRole checks if user role meets the minimum required role
// Role hierarchy: Admin > Manager > User > External
func hasRequiredSystemRole(userRole, requiredRole models.RoleType) bool {
	roleHierarchy := map[models.RoleType]int{
		models.RoleExternal: 0,
		models.RoleUser:     1,
		models.RoleManager:  2,
		models.RoleAdmin:    3,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}
