package router

import (
	chatHandlers "thothix-backend/internal/chat/handlers"
	"thothix-backend/internal/config"
	messageHandlers "thothix-backend/internal/message/handlers"
	"thothix-backend/internal/middleware"
	projectHandlers "thothix-backend/internal/project/handlers"
	sharedHandlers "thothix-backend/internal/shared/handlers"
	sharedMiddleware "thothix-backend/internal/shared/middleware"
	sharedModels "thothix-backend/internal/shared/models"
	userHandlers "thothix-backend/internal/users/handlers"

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
	r.Use(sharedMiddleware.CORS())
	r.Use(sharedMiddleware.Logger())
	r.Use(sharedMiddleware.Recovery())

	// Swagger documentation
	url := ginSwagger.URL("http://localhost:30000/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Health check
	r.GET("/health", sharedHandlers.HealthCheck)

	// Initialize handlers
	authHandler := sharedHandlers.NewAuthHandler(db)
	projectHandler := projectHandlers.NewProjectHandler(db)
	channelHandler := chatHandlers.NewChannelHandler(db)
	messageHandler := messageHandlers.NewMessageHandler(db)
	roleHandler := sharedHandlers.NewRoleHandler(db)
	// API routes
	v1 := r.Group("/api/v1")

	// Public routes (non protette)
	auth := v1.Group("/auth")

	// Webhook di Clerk (middleware + handler pattern)
	auth.POST("/webhooks/clerk",
		sharedMiddleware.ClerkWebhookHandler(cfg.ClerkWebhookSecret),
		authHandler.WebhookHandler,
	)

	// Protected routes con Clerk SDK Auth
	protected := v1.Group("/")
	protected.Use(sharedMiddleware.ClerkAuthSDK(cfg.ClerkSecretKey))
	protected.Use(sharedMiddleware.SetUserContext()) // Add user context for GORM hooks

	// Auth routes (sync with Clerk)
	authProtected := protected.Group("/auth")
	authProtected.POST("/sync", authHandler.SyncUser)
	authProtected.GET("/me", authHandler.GetCurrentUser)
	authProtected.POST("/import-users", middleware.RequireSystemRole(db, sharedModels.RoleAdmin), authHandler.ImportUsers)

	// Users - using the new vertical slice structure
	userHandlers.RegisterUserRoutes(protected, db)

	// Roles management (only admins can manage roles)
	roles := protected.Group("/roles")
	roles.POST("", middleware.RequireSystemRole(db, sharedModels.RoleAdmin), roleHandler.AssignUserRole)
	roles.DELETE("/:roleId", middleware.RequireSystemRole(db, sharedModels.RoleAdmin), roleHandler.RevokeUserRole)

	// Projects
	projects := protected.Group("/projects")
	projects.GET("", projectHandler.GetProjects)
	projects.POST("", middleware.RequirePermission(db, sharedModels.PermissionProjectCreate, nil), projectHandler.CreateProject)
	projects.GET("/:id", middleware.RequireProjectAccess(db), projectHandler.GetProject)
	projects.PUT("/:id", middleware.RequireProjectAccess(db), projectHandler.UpdateProject)
	projects.DELETE("/:id", middleware.RequirePermission(db, sharedModels.PermissionProjectDelete, stringPtr("project")), projectHandler.DeleteProject)
	projects.POST("/:id/members", middleware.RequirePermission(db, sharedModels.PermissionProjectManage, stringPtr("project")), projectHandler.AddMember)
	projects.DELETE("/:id/members/:userId", middleware.RequirePermission(db, sharedModels.PermissionProjectManage, stringPtr("project")), projectHandler.RemoveMember)

	// Channels
	channels := protected.Group("/channels")
	channels.GET("", channelHandler.GetChats)
	channels.POST("", middleware.RequirePermission(db, sharedModels.PermissionChannelCreate, nil), channelHandler.CreateChat)
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

// SetupTestRouter creates a router without authentication middleware for testing
func SetupTestRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Basic middleware for tests
	r.Use(sharedMiddleware.CORS())
	r.Use(sharedMiddleware.Logger())
	r.Use(sharedMiddleware.Recovery())

	// Mock authentication middleware for tests
	r.Use(func(c *gin.Context) {
		// Mock user context - simulates what ClerkAuthSDK middleware would do
		c.Set("clerk_user_id", "test-clerk-user-id")
		c.Set("user_id", "test-user-id")
		c.Next()
	})
	r.Use(sharedMiddleware.SetUserContext()) // Add user context for GORM hooks

	// Health check
	r.GET("/health", sharedHandlers.HealthCheck)

	// Initialize handlers
	authHandler := sharedHandlers.NewAuthHandler(db)
	projectHandler := projectHandlers.NewProjectHandler(db)
	channelHandler := chatHandlers.NewChannelHandler(db)
	messageHandler := messageHandlers.NewMessageHandler(db)

	// API routes
	v1 := r.Group("/api/v1")

	// Auth routes (no authentication required in tests)
	auth := v1.Group("/auth")
	auth.POST("/sync", authHandler.SyncUser)
	auth.GET("/me", authHandler.GetCurrentUser)

	// Users - using the new vertical slice structure (no auth middleware in tests)
	userHandlers.RegisterUserRoutes(v1, db)

	// Projects (simplified for tests)
	projects := v1.Group("/projects")
	projects.GET("", projectHandler.GetProjects)
	projects.POST("", projectHandler.CreateProject)
	projects.GET("/:id", projectHandler.GetProject)
	projects.PUT("/:id", projectHandler.UpdateProject)
	projects.DELETE("/:id", projectHandler.DeleteProject)

	// Channels (simplified for tests)
	channels := v1.Group("/channels")
	channels.GET("", channelHandler.GetChats)
	channels.POST("", channelHandler.CreateChat)
	channels.GET("/:id", channelHandler.GetChat)
	channels.POST("/:id/join", channelHandler.JoinChannel)
	channels.GET("/:id/messages", messageHandler.GetMessages)
	channels.POST("/:id/messages", messageHandler.SendMessage)

	return r
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
