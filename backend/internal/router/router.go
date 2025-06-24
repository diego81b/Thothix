package router

import (
	"thothix-backend/internal/config"
	"thothix-backend/internal/handlers"
	"thothix-backend/internal/middleware"
	"thothix-backend/internal/models"
	"thothix-backend/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middleware globali
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// Swagger documentation
	url := ginSwagger.URL("http://localhost:30000/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Initialize services
	userService := services.NewUserService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(db)
	channelHandler := handlers.NewChannelHandler(db)
	messageHandler := handlers.NewMessageHandler(db)
	roleHandler := handlers.NewRoleHandler(db)
	// API routes
	v1 := r.Group("/api/v1")

	// Public routes (non protette)
	auth := v1.Group("/auth")

	// Webhook di Clerk (middleware + handler pattern)
	auth.POST("/webhooks/clerk",
		middleware.ClerkWebhookHandler(cfg.ClerkWebhookSecret),
		authHandler.WebhookHandler,
	)

	// Protected routes con Clerk SDK Auth
	protected := v1.Group("/")
	protected.Use(middleware.ClerkAuthSDK(cfg.ClerkSecretKey))
	protected.Use(middleware.SetUserContext()) // Add user context for GORM hooks

	// Auth routes (sync with Clerk)
	authProtected := protected.Group("/auth")
	authProtected.POST("/sync", authHandler.SyncUser)
	authProtected.GET("/me", authHandler.GetCurrentUser)
	authProtected.POST("/import-users", middleware.RequireSystemRole(db, models.RoleAdmin), authHandler.ImportUsers)

	// Users
	users := protected.Group("/users")
	users.GET("", userHandler.GetUsers)
	users.PUT("/me", userHandler.UpdateCurrentUser)
	users.GET("/:id", userHandler.GetUser)
	users.PUT("/:id", middleware.RequirePermission(db, models.PermissionUserManage, nil), userHandler.UpdateUser)
	users.GET("/:id/roles", roleHandler.GetUserRoles)

	// Roles management (only admins can manage roles)
	roles := protected.Group("/roles")
	roles.POST("", middleware.RequireSystemRole(db, models.RoleAdmin), roleHandler.AssignUserRole)
	roles.DELETE("/:roleId", middleware.RequireSystemRole(db, models.RoleAdmin), roleHandler.RevokeUserRole)

	// Projects
	projects := protected.Group("/projects")
	projects.GET("", projectHandler.GetProjects)
	projects.POST("", middleware.RequirePermission(db, models.PermissionProjectCreate, nil), projectHandler.CreateProject)
	projects.GET("/:id", middleware.RequireProjectAccess(db), projectHandler.GetProject)
	projects.PUT("/:id", middleware.RequireProjectAccess(db), projectHandler.UpdateProject)
	projects.DELETE("/:id", middleware.RequirePermission(db, models.PermissionProjectDelete, stringPtr("project")), projectHandler.DeleteProject)
	projects.POST("/:id/members", middleware.RequirePermission(db, models.PermissionProjectManage, stringPtr("project")), projectHandler.AddMember)
	projects.DELETE("/:id/members/:userId", middleware.RequirePermission(db, models.PermissionProjectManage, stringPtr("project")), projectHandler.RemoveMember)

	// Channels
	channels := protected.Group("/channels")
	channels.GET("", channelHandler.GetChats)
	channels.POST("", middleware.RequirePermission(db, models.PermissionChannelCreate, nil), channelHandler.CreateChat)
	channels.GET("/:id", middleware.RequireChannelAccess(db), channelHandler.GetChat)
	channels.POST("/:id/join", channelHandler.JoinChannel)
	channels.GET("/:id/messages", middleware.RequireChannelAccess(db), messageHandler.GetMessages)
	channels.POST("/:id/messages", middleware.RequireChannelAccess(db), messageHandler.SendMessage)

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		// TODO: Implement WebSocket handler with Clerk auth
		c.JSON(200, gin.H{"message": "WebSocket endpoint - TODO"})
	})

	return r
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
