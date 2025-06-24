package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"thothix-backend/internal/dto"
	"thothix-backend/internal/middleware"
	"thothix-backend/internal/services"
)

type AuthHandler struct {
	db          *gorm.DB
	userService services.UserServiceInterface
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		db:          db,
		userService: services.NewUserService(db),
	}
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
	clerkUsername, _ := c.Get("clerk_username")
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
		if username, ok := clerkUsername.(*string); ok && username != nil {
			fullName = *username
		}
	}

	// Estrai username
	var username string
	if usernameStr, ok := clerkUsername.(*string); ok && usernameStr != nil {
		username = *usernameStr
	}

	// Prepara i dati per il servizio usando il DTO
	clerkSyncReq := &dto.ClerkUserSyncRequest{
		ClerkID:   clerkUserID.(string),
		Email:     clerkEmail.(string),
		Name:      fullName,
		Username:  username,
		AvatarURL: clerkImageURL.(string),
	}

	// Utilizza il servizio per sincronizzare l'utente
	syncResponse, err := h.userService.SyncUserFromClerk(clerkSyncReq)
	if err != nil {
		log.Printf("Error syncing user: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "sync_failed",
			Message: "Failed to sync user",
		})
		return
	}

	if syncResponse.IsNew {
		c.JSON(http.StatusCreated, syncResponse.User)
	} else {
		c.JSON(http.StatusOK, syncResponse.User)
	}
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

	user, err := h.userService.GetUserByID(clerkUserID.(string))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "user_not_found",
				Message: "User not found",
			})
		} else {
			log.Printf("Error getting user: %v", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "database_error",
				Message: "Database error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
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
	// Get typed webhook event from middleware
	webhookID, _ := middleware.GetWebhookIDFromContext(c)
	log.Printf("Processing Clerk webhook %s", webhookID)

	event, exists := middleware.GetWebhookEventFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook event not found in context"})
		return
	}

	log.Printf("Processing webhook event type: %s", event.Type)

	switch event.Type {
	case "user.created":
		if userData, ok := middleware.GetWebhookUserDataFromContext(c); ok {
			user, err := h.userService.CreateUserFromWebhook(userData)
			if err != nil {
				log.Printf("Error handling user.created webhook %s: %v", webhookID, err)
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "webhook_processing_failed",
					Message: "Failed to process webhook",
				})
				return
			}
			log.Printf("Created user %s from webhook %s", user.ID, webhookID)
		} else {
			log.Printf("Missing user data for user.created webhook %s", webhookID)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "missing_user_data",
				Message: "Missing user data",
			})
			return
		}

	case "user.updated":
		if userData, ok := middleware.GetWebhookUserDataFromContext(c); ok {
			user, err := h.userService.UpdateUserFromWebhook(userData)
			if err != nil {
				log.Printf("Error handling user.updated webhook %s: %v", webhookID, err)
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "webhook_processing_failed",
					Message: "Failed to process webhook",
				})
				return
			}
			log.Printf("Updated user %s from webhook %s", user.ID, webhookID)
		} else {
			log.Printf("Missing user data for user.updated webhook %s", webhookID)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "missing_user_data",
				Message: "Missing user data",
			})
			return
		}

	case "user.deleted":
		if userData, ok := middleware.GetWebhookUserDataFromContext(c); ok {
			if err := h.userService.DeleteUserFromWebhook(userData); err != nil {
				log.Printf("Error handling user.deleted webhook %s: %v", webhookID, err)
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
					Error:   "webhook_processing_failed",
					Message: "Failed to process webhook",
				})
				return
			}
			log.Printf("Deleted user from webhook %s", webhookID)
		} else {
			log.Printf("Missing user data for user.deleted webhook %s", webhookID)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "missing_user_data",
				Message: "Missing user data",
			})
			return
		}

	default:
		log.Printf("Ignoring unhandled webhook event type: %s (ID: %s)", event.Type, webhookID)
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Webhook processed successfully",
	})
}

// ImportUsers importa tutti gli utenti da Clerk al database locale
// @Summary Import all users from Clerk
// @Description Import all users from Clerk to local database (development only)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth/import-users [post]
func (h *AuthHandler) ImportUsers(c *gin.Context) {
	// TODO: Add admin check here
	// Only allow admins to import users

	// This would call Clerk API to get all users and sync them
	// For now, return a placeholder response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "User import functionality - to be implemented. This will call Clerk Users API and sync all users to local DB",
	})
}
