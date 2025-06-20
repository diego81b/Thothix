package middleware

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
)

// ClerkAuthSDK middleware using official Clerk SDK middleware
// This uses the idiomatic WithHeaderAuthorization middleware from clerk/http package
func ClerkAuthSDK(clerkSecretKey string) gin.HandlerFunc {
	// Set the Clerk API key globally (required for all operations)
	clerk.SetKey(clerkSecretKey)

	// Create the official Clerk HTTP middleware
	clerkMiddleware := clerkhttp.WithHeaderAuthorization()

	return func(c *gin.Context) {
		// Convert Gin context to standard HTTP handler format
		wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract session claims that were verified by Clerk middleware
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No valid session claims found"})
				c.Abort()
				return
			}

			// Get user details from Clerk API for additional user information
			ctx := r.Context()
			userDetails, err := user.Get(ctx, claims.Subject)
			if err != nil {
				// Log the error but continue with claims data only
				// This makes the middleware more resilient to Clerk API issues
				c.Set("clerk_user_id", claims.Subject)
				c.Set("user_id", claims.Subject) // For BaseModel hooks
				c.Set("clerk_session_id", claims.SessionID)
				c.Set("clerk_issued_at", claims.IssuedAt)
				if claims.Expiry != nil {
					c.Set("clerk_expires_at", *claims.Expiry)
				}
			} else {
				// Set comprehensive user context when API call succeeds
				c.Set("clerk_user_id", userDetails.ID)
				c.Set("user_id", userDetails.ID) // For BaseModel hooks

				// Handle optional fields safely
				if userDetails.PrimaryEmailAddressID != nil {
					// Find primary email address from the email addresses array
					for _, email := range userDetails.EmailAddresses {
						if email.ID == *userDetails.PrimaryEmailAddressID {
							c.Set("clerk_email", email.EmailAddress)
							break
						}
					}
				}

				if userDetails.Username != nil {
					c.Set("clerk_username", *userDetails.Username)
				}

				if userDetails.FirstName != nil {
					c.Set("clerk_first_name", *userDetails.FirstName)
				}

				if userDetails.LastName != nil {
					c.Set("clerk_last_name", *userDetails.LastName)
				}

				if userDetails.ImageURL != nil {
					c.Set("clerk_image_url", *userDetails.ImageURL)
				}

				// Set additional claims from JWT
				c.Set("clerk_session_id", claims.SessionID)
				c.Set("clerk_issued_at", claims.IssuedAt)
				if claims.Expiry != nil {
					c.Set("clerk_expires_at", *claims.Expiry)
				}
			}

			c.Next()
		})

		// Apply the Clerk middleware, then our wrapped handler
		handler := clerkMiddleware(wrapped)
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// ClerkWebhookHandler for webhook signature verification
// Note: Clerk SDK v2 doesn't expose webhook verification directly
// We implement basic validation here, with TODO for proper Svix verification
func ClerkWebhookHandler(webhookSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Svix headers required for webhook verification
		timestamp := c.GetHeader("svix-timestamp")
		signature := c.GetHeader("svix-signature")
		id := c.GetHeader("svix-id")

		if timestamp == "" || signature == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing required webhook headers (svix-timestamp, svix-signature)",
			})
			return
		}

		// Read request body
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Failed to read request body",
				"details": err.Error(),
			})
			return
		}

		// TODO: Implement proper Svix webhook verification
		// For now, we'll accept the webhook and let the handler process it
		// In production, you should implement proper signature verification
		// using the svix-go library or implement HMAC SHA-256 verification

		// Store webhook data for the handler
		c.Set("webhook_body", body)
		c.Set("webhook_timestamp", timestamp)
		c.Set("webhook_signature", signature)
		c.Set("webhook_id", id)

		c.Next()
	}
}

// GetClerkUserFromContext helper to extract Clerk user data from Gin context
func GetClerkUserFromContext(c *gin.Context) (map[string]interface{}, bool) {
	userData := make(map[string]interface{})

	if userID, exists := c.Get("clerk_user_id"); exists {
		userData["user_id"] = userID
	} else {
		return nil, false
	}

	// Add other optional fields
	if email, exists := c.Get("clerk_email"); exists {
		userData["email"] = email
	}
	if username, exists := c.Get("clerk_username"); exists {
		userData["username"] = username
	}
	if firstName, exists := c.Get("clerk_first_name"); exists {
		userData["first_name"] = firstName
	}
	if lastName, exists := c.Get("clerk_last_name"); exists {
		userData["last_name"] = lastName
	}
	if imageURL, exists := c.Get("clerk_image_url"); exists {
		userData["image_url"] = imageURL
	}
	if sessionID, exists := c.Get("clerk_session_id"); exists {
		userData["session_id"] = sessionID
	}

	return userData, true
}
