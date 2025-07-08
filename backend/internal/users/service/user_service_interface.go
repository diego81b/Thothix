package service

import (
	"thothix-backend/internal/shared/dto"
	sharedMiddleware "thothix-backend/internal/shared/middleware"
	usersDto "thothix-backend/internal/users/dto"
)

// UserServiceInterface defines the contract for core user operations using Response pattern
type UserServiceInterface interface {
	// Core CRUD Operations using Response pattern with lazy evaluation
	GetUserByID(userID string) *usersDto.GetUserResponse
	GetUserByClerkID(clerkID string) *usersDto.GetUserResponse
	GetUsers(req *dto.PaginationRequest) *usersDto.GetUsersResponse
	CreateUser(req *usersDto.CreateUserRequest) *usersDto.CreateUserResponse
	UpdateUser(userID string, req *usersDto.UpdateUserRequest) *usersDto.UpdateUserResponse
	DeleteUser(userID string) *usersDto.DeleteUserResponse
}

// ClerkUserServiceInterface defines the contract for Clerk-specific user operations
type ClerkUserServiceInterface interface {
	// Clerk Integration using Response pattern
	SyncUserFromClerk(req *usersDto.ClerkUserSyncRequest) *usersDto.CreateUserResponse
	ProcessClerkWebhook(userData *sharedMiddleware.UserWebhookData) *usersDto.ClerkSyncUserResponse
}
