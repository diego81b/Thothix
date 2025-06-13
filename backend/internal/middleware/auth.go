package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Imposta le informazioni dell'utente nel context
		c.Set("clerk_user_id", clerkUser.ID)
		c.Set("clerk_email", clerkUser.Email)
		c.Set("clerk_username", clerkUser.Username)

		c.Next()
	}
}

type ClerkUser struct {
	ID       string `json:"id"`
	Email    string `json:"email_addresses[0].email_address"`
	Username string `json:"username"`
}

func verifyClerkToken(token, secretKey string) (*ClerkUser, error) {
	// Crea una richiesta HTTP per verificare il token con Clerk
	req, err := http.NewRequest("GET", "https://api.clerk.dev/v1/me", nil)
	if err != nil {
		return nil, err
	}

	// Aggiungi gli header necessari
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Esegui la richiesta
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("clerk verification failed: %d", resp.StatusCode)
	}

	// Decodifica la risposta
	var clerkUser ClerkUser
	if err := json.NewDecoder(resp.Body).Decode(&clerkUser); err != nil {
		return nil, err
	}

	return &clerkUser, nil
}
