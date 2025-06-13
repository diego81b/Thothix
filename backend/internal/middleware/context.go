package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// contextKey is a type for context keys to avoid collisions
type contextKey string

const userIDKey contextKey = "user_id"

// SetUserContext sets the user ID in the request context
func SetUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from Clerk middleware
		if userID, exists := c.Get("clerk_user_id"); exists {
			// Set user ID in context for GORM hooks
			ctx := context.WithValue(c.Request.Context(), userIDKey, userID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
