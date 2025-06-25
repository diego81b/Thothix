package services

import (
	"thothix-backend/internal/dto"
	"thothix-backend/internal/middleware"
)

// UserServiceInterface defines the contract for core user operations using Response pattern
type UserServiceInterface interface {
	// Core CRUD Operations using Response pattern with lazy evaluation
	GetUserByID(userID string) *dto.GetUserResponse
	GetUserByClerkID(clerkID string) *dto.GetUserResponse
	GetUsers(req *dto.PaginationRequest) *dto.GetUsersResponse
	CreateUser(req *dto.CreateUserRequest) *dto.CreateUserResponse
	UpdateUser(userID string, req *dto.UpdateUserRequest) *dto.UpdateUserResponse
	DeleteUser(userID string) *dto.DeleteUserResponse
}

// ClerkUserServiceInterface defines the contract for Clerk-specific user operations
type ClerkUserServiceInterface interface {
	// Clerk Integration using Response pattern
	SyncUserFromClerk(req *dto.ClerkUserSyncRequest) *dto.CreateUserResponse
	ProcessClerkWebhook(userData *middleware.UserWebhookData) *dto.ClerkSyncUserResponse
}

// Verify that UserService implements both interfaces
var (
	_ UserServiceInterface      = (*UserService)(nil)
	_ ClerkUserServiceInterface = (*UserService)(nil)
)
