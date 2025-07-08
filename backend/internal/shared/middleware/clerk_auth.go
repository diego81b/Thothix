package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// WebhookEvent represents a Clerk webhook event
type WebhookEvent struct {
	Type      string          `json:"type"`
	Object    string          `json:"object"`
	Data      json.RawMessage `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

// UserWebhookData represents user-related webhook data
type UserWebhookData struct {
	ID                    string   `json:"id"`
	Object                string   `json:"object"`
	ExternalID            *string  `json:"external_id"`
	PrimaryEmailAddressID *string  `json:"primary_email_address_id"`
	PrimaryPhoneNumberID  *string  `json:"primary_phone_number_id"`
	Username              *string  `json:"username"`
	FirstName             *string  `json:"first_name"`
	LastName              *string  `json:"last_name"`
	ImageURL              *string  `json:"image_url"`
	CreatedAt             int64    `json:"created_at"`
	UpdatedAt             int64    `json:"updated_at"`
	EmailAddresses        []Email  `json:"email_addresses"`
	PhoneNumbers          []Phone  `json:"phone_numbers"`
	WebURLs               []WebURL `json:"web3_wallets"`
}

// Email represents an email address in webhook data
type Email struct {
	ID           string `json:"id"`
	Object       string `json:"object"`
	EmailAddress string `json:"email_address"`
	Verification struct {
		Status   string `json:"status"`
		Strategy string `json:"strategy"`
	} `json:"verification"`
}

// Phone represents a phone number in webhook data
type Phone struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	PhoneNumber string `json:"phone_number"`
}

// WebURL represents a web3 wallet in webhook data
type WebURL struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	WebURL  string `json:"web3_wallet"`
	Network string `json:"network"`
}

// ClerkWebhookHandler for webhook signature verification with proper typing
// Implements Svix signature verification for production security
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

		// Verify webhook signature using Svix algorithm
		if !verifyWebhookSignature(body, signature, timestamp, webhookSecret) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid webhook signature",
			})
			return
		}

		// Parse webhook event
		var event WebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid webhook payload",
				"details": err.Error(),
			})
			return
		}

		// Validate event timestamp (not older than 5 minutes)
		now := time.Now().Unix()
		if now-event.Timestamp > 300 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Webhook event too old",
			})
			return
		}

		// Store essential webhook data for handlers
		c.Set("webhook_event", event)
		c.Set("webhook_id", id) // Useful for logging and tracing

		// Parse user data for user-related events
		if isUserEvent(event.Type) {
			var userData UserWebhookData
			if err := json.Unmarshal(event.Data, &userData); err == nil {
				c.Set("webhook_user_data", userData)
			}
		}

		c.Next()
	}
}

// verifyWebhookSignature implements Svix signature verification algorithm
func verifyWebhookSignature(payload []byte, signature, timestamp, secret string) bool {
	// Parse timestamp
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	// Create signed payload: timestamp.payload
	signedPayload := fmt.Sprintf("%d.%s", ts, string(payload))

	// Decode the base64-encoded secret
	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return false
	}

	// Compute HMAC SHA-256
	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(signedPayload))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// Parse signatures from header (format: "v1,signature1 v1,signature2")
	signatures := strings.Split(signature, " ")
	for _, sig := range signatures {
		if strings.HasPrefix(sig, "v1,") {
			providedSignature := strings.TrimPrefix(sig, "v1,")
			if hmac.Equal([]byte(expectedSignature), []byte(providedSignature)) {
				return true
			}
		}
	}

	return false
}

// isUserEvent checks if the webhook event is user-related
func isUserEvent(eventType string) bool {
	userEvents := []string{
		"user.created",
		"user.updated",
		"user.deleted",
	}

	for _, userEvent := range userEvents {
		if eventType == userEvent {
			return true
		}
	}
	return false
}

// GetClerkUserFromContext helper to extract Clerk user data from Gin context
func GetClerkUserFromContext(c *gin.Context) (map[string]any, bool) {
	userData := make(map[string]any)

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

// GetWebhookEventFromContext extracts typed webhook event from context
func GetWebhookEventFromContext(c *gin.Context) (*WebhookEvent, bool) {
	if event, exists := c.Get("webhook_event"); exists {
		if webhookEvent, ok := event.(WebhookEvent); ok {
			return &webhookEvent, true
		}
	}
	return nil, false
}

// GetWebhookUserDataFromContext extracts typed user data from webhook context
func GetWebhookUserDataFromContext(c *gin.Context) (*UserWebhookData, bool) {
	if userData, exists := c.Get("webhook_user_data"); exists {
		if webhookUserData, ok := userData.(UserWebhookData); ok {
			return &webhookUserData, true
		}
	}
	return nil, false
}

// GetWebhookIDFromContext extracts webhook ID for logging and tracing
func GetWebhookIDFromContext(c *gin.Context) (string, bool) {
	if id, exists := c.Get("webhook_id"); exists {
		if webhookID, ok := id.(string); ok {
			return webhookID, true
		}
	}
	return "", false
}
