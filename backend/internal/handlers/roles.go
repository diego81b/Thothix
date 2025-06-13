package handlers

import (
	"net/http"
	"thothix-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleHandler struct {
	db *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{db: db}
}

// AssignUserRole godoc
// @Summary Assign role to user
// @Description Assign a role to a user (system, project, or channel specific)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role body AssignRoleRequest true "Role assignment"
// @Success 201 {object} models.UserRole
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /api/v1/roles [post]
func (h *RoleHandler) AssignUserRole(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create role assignment
	userRole := models.UserRole{
		UserID:       req.UserID,
		Role:         req.Role,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
	}

	if err := h.db.Create(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
		return
	}

	c.JSON(http.StatusCreated, userRole)
}

// GetUserRoles godoc
// @Summary Get user roles
// @Description Get all roles assigned to a user
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {array} models.UserRole
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{userId}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("userId")

	var roles []models.UserRole
	if err := h.db.Where("user_id = ?", userID).Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// RevokeUserRole godoc
// @Summary Revoke user role
// @Description Revoke a specific role from a user
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roleId path string true "Role ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/roles/{roleId} [delete]
func (h *RoleHandler) RevokeUserRole(c *gin.Context) {
	roleID := c.Param("roleId")

	if err := h.db.Delete(&models.UserRole{}, "id = ?", roleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke role"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignRoleRequest represents the request body for role assignment
type AssignRoleRequest struct {
	UserID       string          `json:"user_id" binding:"required"`
	Role         models.RoleType `json:"role" binding:"required"`
	ResourceType *string         `json:"resource_type,omitempty"`
	ResourceID   *string         `json:"resource_id,omitempty"`
}
