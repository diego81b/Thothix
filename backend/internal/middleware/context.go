package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// SetUserContext sets the user ID in the GORM context for automatic population
func SetUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from Clerk middleware
		if userID, exists := c.Get("clerk_user_id"); exists {
			// Set user ID in context for GORM hooks
			ctx := context.WithValue(c.Request.Context(), "user_id", userID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
