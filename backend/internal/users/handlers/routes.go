package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"thothix-backend/internal/users/service"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userService := service.NewUserService(db)
	userHandler := NewUserHandler(userService)

	users := router.Group("/users")
	users.GET("/:id", userHandler.GetUserByID)
	users.GET("", userHandler.GetUsers)
	users.POST("", userHandler.CreateUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}
