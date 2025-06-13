package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ClerkAuth middleware per verificare i token di Clerk
func ClerkAuth(clerkSecretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Bearer token format
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		// Verifica il token con Clerk
		clerkUser, err := verifyClerkToken(token, clerkSecretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		// Imposta le informazioni dell'utente nel context
		c.Set("clerk_user_id", clerkUser.ID)
		c.Set("clerk_email", clerkUser.PrimaryEmailAddress.EmailAddress)
		c.Set("clerk_username", clerkUser.Username)
		c.Set("clerk_first_name", clerkUser.FirstName)
		c.Set("clerk_last_name", clerkUser.LastName)
		c.Set("clerk_image_url", clerkUser.ImageURL)
		c.Set("user_id", clerkUser.ID) // Per BaseModel hooks

		c.Next()
	}
}

// Struttura completa per la risposta di Clerk
type ClerkUser struct {
	CreatedAt           int64                  `json:"created_at"`
	UpdatedAt           int64                  `json:"updated_at"`
	LastSignInAt        *int64                 `json:"last_sign_in_at"`
	PrimaryEmailAddress *ClerkEmailAddress     `json:"primary_email_address"`
	PrimaryPhoneNumber  *ClerkPhoneNumber      `json:"primary_phone_number"`
	Username            *string                `json:"username"`
	FirstName           *string                `json:"first_name"`
	LastName            *string                `json:"last_name"`
	EmailAddresses      []ClerkEmailAddress    `json:"email_addresses"`
	PhoneNumbers        []ClerkPhoneNumber     `json:"phone_numbers"`
	PublicMetadata      map[string]interface{} `json:"public_metadata"`
	PrivateMetadata     map[string]interface{} `json:"private_metadata"`
	UnsafeMetadata      map[string]interface{} `json:"unsafe_metadata"`
	ID                  string                 `json:"id"`
	ImageURL            string                 `json:"image_url"`
}

type ClerkEmailAddress struct {
	Verification *ClerkVerification `json:"verification"`
	ID           string             `json:"id"`
	EmailAddress string             `json:"email_address"`
}

type ClerkPhoneNumber struct {
	Verification *ClerkVerification `json:"verification"`
	ID           string             `json:"id"`
	PhoneNumber  string             `json:"phone_number"`
}

type ClerkVerification struct {
	Status   string `json:"status"`
	Strategy string `json:"strategy"`
}

func verifyClerkToken(token, secretKey string) (*ClerkUser, error) {
	// Crea context con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Crea richiesta per ottenere le informazioni dell'utente
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.clerk.com/v1/users/me", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Headers necessari per Clerk
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Esegui la richiesta
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token with Clerk: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("clerk verification failed with status: %d", resp.StatusCode)
	}

	// Decodifica la risposta
	var clerkUser ClerkUser
	if err := json.NewDecoder(resp.Body).Decode(&clerkUser); err != nil {
		return nil, fmt.Errorf("failed to decode Clerk response: %w", err)
	}

	// Validazione dati essenziali
	if clerkUser.ID == "" {
		return nil, fmt.Errorf("invalid user data: missing ID")
	}

	if clerkUser.PrimaryEmailAddress == nil || clerkUser.PrimaryEmailAddress.EmailAddress == "" {
		return nil, fmt.Errorf("invalid user data: missing primary email")
	}

	return &clerkUser, nil
}
