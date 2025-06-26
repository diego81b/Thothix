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
	db               *gorm.DB
	userService      services.UserServiceInterface
	clerkUserService services.ClerkUserServiceInterface
	userServiceImpl  *services.UserService // For legacy webhook methods
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	userServiceImpl := services.NewUserService(db)
	return &AuthHandler{
		db:               db,
		userService:      userServiceImpl,
		clerkUserService: userServiceImpl,
		userServiceImpl:  userServiceImpl,
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
// @Failure 400 {object} dto.ErrorViewModel
// @Failure 500 {object} dto.ErrorViewModel
// @Router /api/v1/auth/sync [post]
func (h *AuthHandler) SyncUser(c *gin.Context) {
	ctx := WrapContext(c)

	clerkUserID, exists := c.Get("clerk_user_id")
	if !exists {
		ctx.BadRequestErrorResponse("Clerk user ID not found")
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
	output := h.clerkUserService.SyncUserFromClerk(clerkSyncReq)

	output.Match(
		// Exception
		func(exception error) interface{} {
			ctx.SystemErrorResponse(exception, "System error syncing user")
			return nil
		},
		// Success
		func(success *dto.UserResponse) interface{} {
			ctx.SuccessResponse(success)
			return nil
		},
		// Failure
		func(errors []dto.Error) interface{} {
			ctx.ValidationErrorResponse(errors, "Validation errors syncing user")
			return nil
		},
	)
}

// GetCurrentUser ottiene l'utente corrente autenticato
// @Summary Get current user
// @Description Get the current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} dto.ErrorViewModel
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	ctx := WrapContext(c)

	clerkUserID, exists := c.Get("clerk_user_id")
	if !exists {
		ctx.BadRequestErrorResponse("Clerk user ID not found")
		return
	}

	output := h.userService.GetUserByID(clerkUserID.(string))

	output.Match(
		// Exception
		func(exception error) interface{} {
			ctx.SystemErrorResponse(exception, "System error getting user")
			return nil
		},
		// Success
		func(success *dto.UserResponse) interface{} {
			ctx.SuccessResponse(success)
			return nil
		},
		// Failure
		func(errors []dto.Error) interface{} {
			// Check if it's a "not found" error
			for _, err := range errors {
				if err.Message == "User not found" {
					ctx.NotFoundErrorResponse("User", clerkUserID.(string))
					return nil
				}
			}
			ctx.ValidationErrorResponse(errors, "Validation errors getting user")
			return nil
		},
	)
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
			response := h.clerkUserService.ProcessClerkWebhook(userData)
			response.Match(
				func(err error) interface{} {
					c.JSON(http.StatusInternalServerError, dto.LoggedSystemErrorResponse(err, "Error handling user.created webhook %s", webhookID))
					return nil
				},
				func(syncResponse *dto.ClerkUserSyncResponse) interface{} {
					log.Printf("Created user %s from webhook %s", syncResponse.User.ID, webhookID)
					return nil
				},
				func(errors []dto.Error) interface{} {
					c.JSON(http.StatusBadRequest, dto.LoggedValidationErrorResponse(errors, "Validation error handling user.created webhook %s", webhookID))
					return nil
				},
			)
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
			response := h.clerkUserService.ProcessClerkWebhook(userData)
			response.Match(
				func(err error) interface{} {
					c.JSON(http.StatusInternalServerError, dto.LoggedSystemErrorResponse(err, "Error handling user.updated webhook %s", webhookID))
					return nil
				},
				func(syncResponse *dto.ClerkUserSyncResponse) interface{} {
					log.Printf("Updated user %s from webhook %s", syncResponse.User.ID, webhookID)
					return nil
				},
				func(errors []dto.Error) interface{} {
					c.JSON(http.StatusBadRequest, dto.LoggedValidationErrorResponse(errors, "Validation error handling user.updated webhook %s", webhookID))
					return nil
				},
			)
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
			// First find the user by Clerk ID to get internal ID
			getUserResponse := h.userService.GetUserByClerkID(userData.ID)
			var userID string

			getUserResponse.Match(
				func(err error) interface{} {
					c.JSON(http.StatusInternalServerError, dto.LoggedSystemErrorResponse(err, "Error finding user for deletion webhook %s", webhookID))
					return nil
				},
				func(user *dto.UserResponse) interface{} {
					userID = user.ID
					return nil
				},
				func(errors []dto.Error) interface{} {
					log.Printf("User not found for deletion webhook %s: %v", webhookID, errors)
					// User already doesn't exist, consider it successful
					return nil
				},
			)

			if userID != "" {
				deleteResponse := h.userService.DeleteUser(userID)
				deleteResponse.Match(
					func(err error) interface{} {
						c.JSON(http.StatusInternalServerError, dto.LoggedSystemErrorResponse(err, "Error deleting user from webhook %s", webhookID))
						return nil
					},
					func(message string) interface{} {
						log.Printf("Deleted user from webhook %s", webhookID)
						return nil
					},
					func(errors []dto.Error) interface{} {
						c.JSON(http.StatusBadRequest, dto.LoggedValidationErrorResponse(errors, "Validation error deleting user from webhook %s", webhookID))
						return nil
					},
				)
			}
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

	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Webhook processed successfully",
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
	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User import functionality - to be implemented. This will call Clerk Users API and sync all users to local DB",
	})
}
