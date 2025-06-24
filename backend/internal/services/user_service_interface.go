package services

import (
	"thothix-backend/internal/dto"
	"thothix-backend/internal/middleware"
)

// UserServiceInterface defines the contract for user service operations
type UserServiceInterface interface {
	// Clerk Integration
	SyncUserFromClerk(req *dto.ClerkUserSyncRequest) (*dto.ClerkUserSyncResponse, error)
	CreateUserFromWebhook(userData *middleware.UserWebhookData) (*dto.UserResponse, error)
	UpdateUserFromWebhook(userData *middleware.UserWebhookData) (*dto.UserResponse, error)
	DeleteUserFromWebhook(userData *middleware.UserWebhookData) error

	// CRUD Operations
	GetUserByID(userID string) (*dto.UserResponse, error)
	GetUserByClerkID(clerkID string) (*dto.UserResponse, error)
	GetUsers(req *dto.GetUsersRequest) (*dto.UserListResponse, error)
	UpdateUser(userID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
}

// Verify that UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)
