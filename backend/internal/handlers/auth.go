package handlers

import (
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
	clerkUsername, _ := c.Get("clerk_username")

	// Cerca l'utente esistente
	var user models.User
	result := h.db.Where("id = ?", clerkUserID).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// Crea nuovo utente
		user = models.User{
			BaseModel: models.BaseModel{
				ID: clerkUserID.(string),
			},
			Email:     clerkEmail.(string),
			Name:      clerkUsername.(string), // Uso username come name per ora
			AvatarURL: "",                     // Da aggiornare successivamente
		}

		if err := h.db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, user.ToResponse())
		return
	}

	// Aggiorna l'utente esistente se necessario
	updated := false
	if user.Email != clerkEmail.(string) {
		user.Email = clerkEmail.(string)
		updated = true
	}
	if user.Name != clerkUsername.(string) {
		user.Name = clerkUsername.(string)
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
