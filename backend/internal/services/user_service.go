package services

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"thothix-backend/internal/middleware"
	"thothix-backend/internal/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
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
func (s *UserService) SyncUserFromClerk(clerkData ClerkUserData) (*models.User, bool, error) {
	var user models.User
	result := s.db.Where("id = ?", clerkData.ClerkID).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// Crea nuovo utente
		user = models.User{
			BaseModel: models.BaseModel{
				ID: clerkData.ClerkID,
			},
			ClerkID:    clerkData.ClerkID,
			Email:      clerkData.Email,
			Name:       clerkData.Name,
			Username:   clerkData.Username,
			AvatarURL:  clerkData.AvatarURL,
			SystemRole: models.RoleUser,
			LastSync:   time.Now(),
		}

		if err := s.db.Create(&user).Error; err != nil {
			return nil, false, err
		}

		return &user, true, nil
	}

	if result.Error != nil {
		return nil, false, result.Error
	}

	// Aggiorna l'utente esistente se necessario
	updated := false

	if user.Email != clerkData.Email {
		user.Email = clerkData.Email
		updated = true
	}

	if user.Name != clerkData.Name {
		user.Name = clerkData.Name
		updated = true
	}

	if user.Username != clerkData.Username {
		user.Username = clerkData.Username
		updated = true
	}

	if user.AvatarURL != clerkData.AvatarURL {
		user.AvatarURL = clerkData.AvatarURL
		updated = true
	}

	if updated {
		user.LastSync = time.Now()
		if err := s.db.Save(&user).Error; err != nil {
			return nil, false, err
		}
	}

	return &user, false, nil // false = updated or no changes
}

// CreateUserFromWebhook crea un utente dal webhook di Clerk
func (s *UserService) CreateUserFromWebhook(userData *middleware.UserWebhookData) error {
	// Estrai email primaria
	email := s.extractPrimaryEmail(userData)

	// Costruisci il nome completo
	name := s.buildFullName(userData)

	// Estrai avatar URL
	var avatarURL string
	if userData.ImageURL != nil {
		avatarURL = *userData.ImageURL
	}

	// Estrai username
	var username string
	if userData.Username != nil {
		username = *userData.Username
	}

	// Crea utente nel database
	user := models.User{
		BaseModel: models.BaseModel{
			ID: userData.ID,
		},
		ClerkID:    userData.ID,
		Email:      email,
		Name:       name,
		Username:   username,
		AvatarURL:  avatarURL,
		SystemRole: models.RoleUser,
		LastSync:   time.Now(),
	}

	return s.db.Create(&user).Error
}

// UpdateUserFromWebhook aggiorna un utente dal webhook di Clerk
func (s *UserService) UpdateUserFromWebhook(userData *middleware.UserWebhookData) error {
	var user models.User
	if err := s.db.Where("clerk_id = ?", userData.ID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Se l'utente non esiste, crealo
			return s.CreateUserFromWebhook(userData)
		}
		return err
	}

	// Aggiorna i campi
	updated := false

	// Estrai email primaria
	email := s.extractPrimaryEmail(userData)
	if user.Email != email {
		user.Email = email
		updated = true
	}

	// Costruisci il nome completo
	name := s.buildFullName(userData)
	if user.Name != name {
		user.Name = name
		updated = true
	}

	// Estrai username
	var username string
	if userData.Username != nil {
		username = *userData.Username
	}
	if user.Username != username {
		user.Username = username
		updated = true
	}

	// Estrai avatar URL
	var avatarURL string
	if userData.ImageURL != nil {
		avatarURL = *userData.ImageURL
	}
	if user.AvatarURL != avatarURL {
		user.AvatarURL = avatarURL
		updated = true
	}

	if updated {
		user.LastSync = time.Now()
		return s.db.Save(&user).Error
	}

	return nil
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
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByClerkID ottiene un utente per Clerk ID
func (s *UserService) GetUserByClerkID(clerkID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("clerk_id = ?", clerkID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers ottiene tutti gli utenti con paginazione
func (s *UserService) GetUsers(offset, limit int) ([]models.User, error) {
	var users []models.User
	if err := s.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser aggiorna un utente esistente
func (s *UserService) UpdateUser(userID string, updates map[string]interface{}) error {
	result := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
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
