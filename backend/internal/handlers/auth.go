package handlers

import (
	"fmt"
	"log"
	"net/http"

	"thothix-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// SyncUser sincronizza l'utente da Clerk con il database locale
// @Summary Sync user from Clerk
// @Description Synchronize user data from Clerk to local database
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/sync [post]
func (h *AuthHandler) SyncUser(c *gin.Context) {
	clerkUserID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Clerk user ID not found"})
		return
	}

	clerkEmail, _ := c.Get("clerk_email")
	clerkFirstName, _ := c.Get("clerk_first_name")
	clerkLastName, _ := c.Get("clerk_last_name")
	clerkImageURL, _ := c.Get("clerk_image_url")

	// Costruisci il nome completo
	var fullName string
	if firstName, ok := clerkFirstName.(*string); ok && firstName != nil {
		fullName = *firstName
		if lastName, ok := clerkLastName.(*string); ok && lastName != nil {
			fullName += " " + *lastName
		}
	}
	if fullName == "" {
		if username, exists := c.Get("clerk_username"); exists {
			if usernameStr, ok := username.(*string); ok && usernameStr != nil {
				fullName = *usernameStr
			}
		}
	}

	// Cerca l'utente esistente
	var user models.User
	result := h.db.Where("id = ?", clerkUserID).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// Crea nuovo utente
		user = models.User{
			BaseModel: models.BaseModel{
				ID: clerkUserID.(string),
			},
			Email:      clerkEmail.(string),
			Name:       fullName,
			AvatarURL:  clerkImageURL.(string),
			SystemRole: models.RoleUser, // Ruolo di default
		}

		if err := h.db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, user.ToResponse())
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Aggiorna l'utente esistente se necessario
	updated := false

	if user.Email != clerkEmail.(string) {
		user.Email = clerkEmail.(string)
		updated = true
	}

	if user.Name != fullName {
		user.Name = fullName
		updated = true
	}

	if user.AvatarURL != clerkImageURL.(string) {
		user.AvatarURL = clerkImageURL.(string)
		updated = true
	}

	if updated {
		if err := h.db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// GetCurrentUser ottiene l'utente corrente autenticato
// @Summary Get current user
// @Description Get the current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	clerkUserID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Clerk user ID not found"})
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", clerkUserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// WebhookHandler gestisce i webhook di Clerk per sincronizzazione automatica
// @Summary Handle Clerk webhooks
// @Description Handle Clerk webhooks for automatic user synchronization
// @Tags auth
// @Accept json
// @Produce json
// @Param webhook body map[string]interface{} true "Clerk webhook payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/webhooks/clerk [post]
func (h *AuthHandler) WebhookHandler(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	eventType, ok := payload["type"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing event type"})
		return
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing event data"})
		return
	}

	switch eventType {
	case "user.created":
		if err := h.handleUserCreated(data); err != nil {
			log.Printf("Error handling user.created webhook: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
			return
		}
	case "user.updated":
		if err := h.handleUserUpdated(data); err != nil {
			log.Printf("Error handling user.updated webhook: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
			return
		}
	case "user.deleted":
		if err := h.handleUserDeleted(data); err != nil {
			log.Printf("Error handling user.deleted webhook: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process webhook"})
			return
		}
	default:
		// Ignora eventi non gestiti
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed"})
}

func (h *AuthHandler) handleUserCreated(data map[string]interface{}) error {
	userID, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("missing user ID")
	}

	// Estrai email primaria
	var email string
	if emailAddresses, ok := data["email_addresses"].([]interface{}); ok {
		for _, addr := range emailAddresses {
			if addrMap, ok := addr.(map[string]interface{}); ok {
				if emailAddr, ok := addrMap["email_address"].(string); ok {
					email = emailAddr
					break
				}
			}
		}
	}

	// Estrai nome
	var name string
	if firstName, ok := data["first_name"].(string); ok {
		name = firstName
		if lastName, ok := data["last_name"].(string); ok {
			name += " " + lastName
		}
	}
	if name == "" {
		if username, ok := data["username"].(string); ok {
			name = username
		}
	}

	// Estrai avatar
	avatarURL, _ := data["image_url"].(string)

	// Crea utente nel database
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userID,
		},
		Email:      email,
		Name:       name,
		AvatarURL:  avatarURL,
		SystemRole: models.RoleUser,
	}

	return h.db.Create(&user).Error
}

func (h *AuthHandler) handleUserUpdated(data map[string]interface{}) error {
	userID, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("missing user ID")
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Se l'utente non esiste, crealo
			return h.handleUserCreated(data)
		}
		return err
	}

	// Aggiorna i campi
	updated := false

	// Estrai email primaria
	if emailAddresses, ok := data["email_addresses"].([]interface{}); ok {
		for _, addr := range emailAddresses {
			if addrMap, ok := addr.(map[string]interface{}); ok {
				if emailAddr, ok := addrMap["email_address"].(string); ok {
					if user.Email != emailAddr {
						user.Email = emailAddr
						updated = true
					}
					break
				}
			}
		}
	}

	// Estrai nome
	var name string
	if firstName, ok := data["first_name"].(string); ok {
		name = firstName
		if lastName, ok := data["last_name"].(string); ok {
			name += " " + lastName
		}
	}
	if name == "" {
		if username, ok := data["username"].(string); ok {
			name = username
		}
	}
	if user.Name != name {
		user.Name = name
		updated = true
	}

	// Estrai avatar
	if avatarURL, ok := data["image_url"].(string); ok {
		if user.AvatarURL != avatarURL {
			user.AvatarURL = avatarURL
			updated = true
		}
	}

	if updated {
		return h.db.Save(&user).Error
	}

	return nil
}

func (h *AuthHandler) handleUserDeleted(data map[string]interface{}) error {
	userID, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("missing user ID")
	}

	// Invece di eliminare l'utente, potremmo marcarlo come inattivo
	// o gestire la cancellazione in modo diverso secondo le business rules
	return h.db.Where("id = ?", userID).Delete(&models.User{}).Error
}
