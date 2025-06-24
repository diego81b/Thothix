package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"thothix-backend/internal/dto"
	"thothix-backend/internal/mappers"
	"thothix-backend/internal/middleware"
	"thothix-backend/internal/models"
)

type UserService struct {
	db     *gorm.DB
	mapper *mappers.UserMapper
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db:     db,
		mapper: mappers.NewUserMapper(),
	}
}

// ClerkUserData rappresenta i dati utente estratti dal context di Clerk
type ClerkUserData struct {
	ClerkID   string
	Email     string
	Name      string
	Username  string
	AvatarURL string
}

// SyncUserFromClerk sincronizza un utente da Clerk al database locale
func (s *UserService) SyncUserFromClerk(req *dto.ClerkUserSyncRequest) (*dto.ClerkUserSyncResponse, error) {
	var user models.User
	err := s.db.Where("clerk_id = ?", req.ClerkID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// Crea nuovo utente
		user = *s.mapper.ClerkSyncRequestToModel(req)
		// Non impostare manualmente l'ID - lascia che sia generato automaticamente
		user.LastSync = time.Now()

		if err := s.db.Create(&user).Error; err != nil {
			return nil, err
		}

		return s.mapper.CreateSyncResponse(&user, true, "User created successfully"), nil
	}

	if err != nil {
		return nil, err
	}

	// Aggiorna l'utente esistente se necessario
	updated := false

	if user.Email != req.Email {
		user.Email = req.Email
		updated = true
	}

	if user.Name != req.Name {
		user.Name = req.Name
		updated = true
	}

	if user.Username != req.Username {
		user.Username = req.Username
		updated = true
	}

	if user.AvatarURL != req.AvatarURL {
		user.AvatarURL = req.AvatarURL
		updated = true
	}

	if updated {
		user.LastSync = time.Now()
		if err := s.db.Save(&user).Error; err != nil {
			return nil, err
		}
		return s.mapper.CreateSyncResponse(&user, false, "User updated successfully"), nil
	}

	return s.mapper.CreateSyncResponse(&user, false, "User already up to date"), nil
}

// CreateUserFromWebhook crea un utente dal webhook di Clerk
func (s *UserService) CreateUserFromWebhook(userData *middleware.UserWebhookData) (*dto.UserResponse, error) {
	// Converti webhook data in ClerkUserSyncRequest
	syncReq := &dto.ClerkUserSyncRequest{
		ClerkID:   userData.ID,
		Email:     s.extractPrimaryEmail(userData),
		Name:      s.buildFullName(userData),
		AvatarURL: s.extractAvatarURL(userData),
		Username:  s.extractUsername(userData),
	}

	// Usa il metodo sync per creare l'utente
	syncResponse, err := s.SyncUserFromClerk(syncReq)
	if err != nil {
		return nil, err
	}

	return &syncResponse.User, nil
}

// UpdateUserFromWebhook aggiorna un utente dal webhook di Clerk
func (s *UserService) UpdateUserFromWebhook(userData *middleware.UserWebhookData) (*dto.UserResponse, error) {
	// Converti webhook data in ClerkUserSyncRequest
	syncReq := &dto.ClerkUserSyncRequest{
		ClerkID:   userData.ID,
		Email:     s.extractPrimaryEmail(userData),
		Name:      s.buildFullName(userData),
		AvatarURL: s.extractAvatarURL(userData),
		Username:  s.extractUsername(userData),
	}

	// Usa il metodo sync per aggiornare l'utente
	syncResponse, err := s.SyncUserFromClerk(syncReq)
	if err != nil {
		return nil, err
	}

	return &syncResponse.User, nil
}

// DeleteUserFromWebhook gestisce la cancellazione di un utente dal webhook di Clerk
func (s *UserService) DeleteUserFromWebhook(userData *middleware.UserWebhookData) error {
	// Strategia: soft delete o mark as inactive
	// In base alle business rules, potresti voler:
	// 1. Soft delete (GORM built-in)
	// 2. Mark as inactive
	// 3. Hard delete (non raccomandato)

	// Per ora implementiamo soft delete
	return s.db.Where("clerk_id = ?", userData.ID).Delete(&models.User{}).Error
}

// GetUserByID ottiene un utente per ID
func (s *UserService) GetUserByID(userID string) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return s.mapper.ModelToResponse(&user), nil
}

// GetUserByClerkID ottiene un utente per Clerk ID
func (s *UserService) GetUserByClerkID(clerkID string) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Where("clerk_id = ?", clerkID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return s.mapper.ModelToResponse(&user), nil
}

// GetUsers ottiene tutti gli utenti con paginazione
func (s *UserService) GetUsers(req *dto.GetUsersRequest) (*dto.UserListResponse, error) {
	var users []models.User
	var total int64

	// Count total users
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PerPage

	// Get users with pagination
	if err := s.db.Offset(offset).Limit(req.PerPage).Find(&users).Error; err != nil {
		return nil, err
	}

	return s.mapper.ModelsToListResponse(users, total, req.Page, req.PerPage), nil
}

// UpdateUser aggiorna un utente esistente
func (s *UserService) UpdateUser(userID string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	updates := s.mapper.UpdateRequestToMap(req)
	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	result := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	// Retrieve updated user
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.mapper.ModelToResponse(&user), nil
}

// Helper methods

func (s *UserService) extractPrimaryEmail(userData *middleware.UserWebhookData) string {
	if userData.PrimaryEmailAddressID != nil {
		for _, emailAddr := range userData.EmailAddresses {
			if emailAddr.ID == *userData.PrimaryEmailAddressID {
				return emailAddr.EmailAddress
			}
		}
	}
	// Fallback alla prima email se non c'è primary
	if len(userData.EmailAddresses) > 0 {
		return userData.EmailAddresses[0].EmailAddress
	}
	return ""
}

func (s *UserService) buildFullName(userData *middleware.UserWebhookData) string {
	var name string
	if userData.FirstName != nil {
		name = *userData.FirstName
		if userData.LastName != nil {
			name += " " + *userData.LastName
		}
	}
	// Fallback al username se non c'è nome
	if name == "" && userData.Username != nil {
		name = *userData.Username
	}
	return name
}

func (s *UserService) extractUsername(userData *middleware.UserWebhookData) string {
	if userData.Username != nil {
		return *userData.Username
	}
	return ""
}

func (s *UserService) extractAvatarURL(userData *middleware.UserWebhookData) string {
	if userData.ImageURL != nil {
		return *userData.ImageURL
	}
	return ""
}
