package handlers

import (
	"net/http"
	"strconv"

	chatDomain "thothix-backend/internal/chat/domain"
	messageDomain "thothix-backend/internal/message/domain"
	messageDto "thothix-backend/internal/message/dto"
	sharedModels "thothix-backend/internal/shared/models"
	usersDomain "thothix-backend/internal/users/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageHandler struct {
	db *gorm.DB
}

func NewMessageHandler(db *gorm.DB) *MessageHandler {
	return &MessageHandler{db: db}
}

// GetMessages godoc
// @Summary Get messages for a channel
// @Description Get all messages for a specific channel with pagination
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Messages per page" default(50)
// @Success 200 {object} MessageListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/channels/{id}/messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	channelID := c.Param("id")

	// Parse pagination parameters
	page := 1
	limit := 50
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// Check if user has access to this channel (already done by middleware, but double-check)
	resourceType := "channel"
	if !sharedModels.HasUserPermission(h.db, userID.(string), sharedModels.PermissionChannelRead, &resourceType, &channelID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to channel"})
		return
	}

	// Get messages with pagination
	offset := (page - 1) * limit
	var messages []messageDomain.Message
	var total int64

	// Count total messages
	h.db.Model(&messageDomain.Message{}).Where("channel_id = ?", channelID).Count(&total)

	// Get paginated messages with user preloading
	if err := h.db.Preload("Sender").Where("channel_id = ?", channelID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	response := MessageListResponse{
		Messages: messages,
		Page:     page,
		Limit:    limit,
		Total:    total,
		Pages:    (total + int64(limit) - 1) / int64(limit),
	}

	c.JSON(http.StatusOK, response)
}

// SendMessage godoc
// @Summary Send a message
// @Description Send a message to a channel
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel ID"
// @Param message body messageDto.MessageCreateRequest true "Message data"
// @Success 201 {object} models.Message
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/channels/{id}/messages [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	channelID := c.Param("id")

	var req messageDto.MessageCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has permission to send messages in this channel
	resourceType := "channel"
	if !sharedModels.HasUserPermission(h.db, userID.(string), sharedModels.PermissionMessageCreate, &resourceType, &channelID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot send messages to this channel"})
		return
	}

	// Verify channel exists
	var channel chatDomain.Channel
	if err := h.db.Where("id = ?", channelID).First(&channel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	// Create message
	message := messageDomain.Message{
		Content:   req.Content,
		ChannelID: &channelID,
		SenderID:  userID.(string),
	}

	if err := h.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	// Load sender relation for response
	h.db.Preload("Sender").First(&message, message.ID)

	c.JSON(http.StatusCreated, message)
}

// CreateDirectMessage godoc
// @Summary Create/Send direct message
// @Description Create a direct message conversation or send a message to existing DM
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body DirectMessageRequest true "Direct message data"
// @Success 201 {object} models.Message
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/messages/direct [post]
func (h *MessageHandler) CreateDirectMessage(c *gin.Context) {
	userID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req DirectMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has permission to create direct messages
	if !sharedModels.HasUserPermission(h.db, userID.(string), sharedModels.PermissionDMCreate, nil, nil) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot create direct messages"})
		return
	}

	// Verify recipient exists
	var recipient usersDomain.User
	if err := h.db.Where("id = ?", req.RecipientID).First(&recipient).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Recipient not found"})
		return
	}

	// Create direct message
	message := messageDomain.Message{
		Content:    req.Content,
		SenderID:   userID.(string),
		ReceiverID: &req.RecipientID,
	}

	if err := h.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send direct message"})
		return
	}

	// Load user relation for response
	h.db.Preload("Sender").Preload("Receiver").First(&message, message.ID)

	c.JSON(http.StatusCreated, message)
}

// DirectMessageRequest represents the request body for direct messages
type DirectMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	RecipientID string `json:"recipient_id" binding:"required"`
}

// MessageListResponse represents the response for message listing
type MessageListResponse struct {
	Messages []messageDomain.Message `json:"messages"`
	Page     int                     `json:"page"`
	Limit    int                     `json:"limit"`
	Total    int64                   `json:"total"`
	Pages    int64                   `json:"pages"`
}
