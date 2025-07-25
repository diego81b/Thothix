package handlers

import (
	"log"
	"net/http"

	chatDomain "thothix-backend/internal/chat/domain"
	chatDto "thothix-backend/internal/chat/dto"
	projectDomain "thothix-backend/internal/project/domain"
	sharedModels "thothix-backend/internal/shared/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChannelHandler struct {
	db *gorm.DB
}

func NewChannelHandler(db *gorm.DB) *ChannelHandler {
	return &ChannelHandler{db: db}
}

// GetChats godoc
// @Summary Get all channels for user
// @Description Get a list of all channels accessible to the authenticated user
// @Tags channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} chatDomain.Channel
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/channels [get]
func (h *ChannelHandler) GetChats(c *gin.Context) {
	userID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user's system role
	userRole, err := sharedModels.GetUserRole(h.db, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user role"})
		return
	}

	var channels []chatDomain.Channel

	// Admins and managers can see all channels
	switch userRole {
	case sharedModels.RoleAdmin, sharedModels.RoleManager:
		if err := h.db.Preload("Project").Find(&channels).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channels"})
			return
		}
	case sharedModels.RoleExternal:
		// External users can only see public channels (channels without members)
		query := `
			SELECT c.* FROM channels c
			LEFT JOIN channel_members cm ON c.id = cm.channel_id
			WHERE cm.channel_id IS NULL
		`
		if err := h.db.Preload("Project").Raw(query).Scan(&channels).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channels"})
			return
		}
	default:
		// Regular users can see public channels and channels they're members of
		query := `
			SELECT DISTINCT c.* FROM channels c
			LEFT JOIN channel_members cm1 ON c.id = cm1.channel_id AND cm1.user_id = ?
			LEFT JOIN channel_members cm2 ON c.id = cm2.channel_id
			WHERE cm2.channel_id IS NULL OR cm1.user_id IS NOT NULL
		`
		if err := h.db.Preload("Project").Raw(query, userID).Scan(&channels).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get channels"})
			return
		}
	}

	// Load IsPrivate field for each channel
	for i := range channels {
		if err := channels[i].LoadIsPrivate(h.db); err != nil {
			log.Printf("Error loading IsPrivate for channel %s: %v", channels[i].ID, err)
		}
	}

	c.JSON(http.StatusOK, channels)
}

// CreateChat godoc
// @Summary Create a new channel
// @Description Create a new channel for a project
// @Tags channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param channel body chatDto.ChannelCreateRequest true "Channel data"
// @Success 201 {object} chatDomain.Channel
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/channels [post]
func (h *ChannelHandler) CreateChat(c *gin.Context) {
	userID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req chatDto.ChannelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has permission to create channels
	if !sharedModels.HasUserPermission(h.db, userID.(string), sharedModels.PermissionChannelCreate, nil, nil) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions to create channels"})
		return
	}

	// Verify project exists and user has access
	var project projectDomain.Project
	if err := h.db.Where("id = ?", req.ProjectID).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project not found"})
		return
	}

	resourceType := "project"
	if !sharedModels.HasUserPermission(h.db, userID.(string), sharedModels.PermissionProjectRead, &resourceType, &req.ProjectID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to project"})
		return
	}

	// Create channel
	channel := chatDomain.Channel{
		Name:      req.Name,
		ProjectID: req.ProjectID,
	}

	if err := h.db.Create(&channel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
		return
	}

	// Load project relation and IsPrivate field for response
	h.db.Preload("Project").First(&channel, channel.ID)
	if err := channel.LoadIsPrivate(h.db); err != nil {
		log.Printf("Error loading IsPrivate for channel %s: %v", channel.ID, err)
	}

	c.JSON(http.StatusCreated, channel)
}

// GetChat godoc
// @Summary Get channel by ID
// @Description Get a single channel by its ID
// @Tags channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel ID"
// @Success 200 {object} chatDomain.Channel
// @Failure 404 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/channels/{id} [get]
func (h *ChannelHandler) GetChat(c *gin.Context) {
	userID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	channelID := c.Param("id")

	// Check if user has access to this channel
	resourceType := "channel"
	if !sharedModels.HasUserPermission(h.db, userID.(string), sharedModels.PermissionChannelRead, &resourceType, &channelID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to channel"})
		return
	}

	var channel chatDomain.Channel
	if err := h.db.Preload("Project").Where("id = ?", channelID).First(&channel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	c.JSON(http.StatusOK, channel)
}

// JoinChannel godoc
// @Summary Join a channel
// @Description Join a public channel or accept invite to private channel
// @Tags channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel ID"
// @Success 201 {object} chatDomain.ChannelMember
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/channels/{id}/join [post]
func (h *ChannelHandler) JoinChannel(c *gin.Context) {
	userID, exists := c.Get("clerk_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	channelID := c.Param("id")

	// Get channel info
	var channel chatDomain.Channel
	if err := h.db.Where("id = ?", channelID).First(&channel).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	// Get user role
	userRole, err := sharedModels.GetUserRole(h.db, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user role"})
		return
	}

	// Check if user can join this channel
	switch {
	case channel.IsPrivate:
		// Only admins/managers can join private channels without invitation
		if userRole != sharedModels.RoleAdmin && userRole != sharedModels.RoleManager {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot join private channel without invitation"})
			return
		}
	case !channel.IsPrivate && userRole == sharedModels.RoleExternal:
		// External users can join public channels
	case userRole == sharedModels.RoleUser:
		// Regular users can join any public channel if they have project access
		resourceType := "project"
		if !sharedModels.HasUserPermission(h.db, userID.(string), sharedModels.PermissionProjectRead, &resourceType, &channel.ProjectID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to project"})
			return
		}
	}

	// Check if already a member
	var existingMember chatDomain.ChannelMember
	if err := h.db.Where("channel_id = ? AND user_id = ?", channelID, userID).First(&existingMember).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already a member of this channel"})
		return
	}

	// Create membership
	member := chatDomain.ChannelMember{
		ChannelID: channelID,
		UserID:    userID.(string),
	}

	if err := h.db.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join channel"})
		return
	}

	c.JSON(http.StatusCreated, member)
}
