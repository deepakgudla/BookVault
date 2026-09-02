package repository

import (
	"time"

	"github.com/deepakgudla/bookvault/internal/models"
	"gorm.io/gorm"
)

// UserRepository persists users and refresh tokens.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a user repository backed by the supplied database.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByEmail returns the user with the given email address.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID returns the user with the given ID.
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmailAndActive returns a user matching an email address and active state.
func (r *UserRepository) GetByEmailAndActive(email string, isActive bool) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND is_active = ?", email, isActive).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create persists a new user.
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update persists changes to an existing user.
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete removes a user by ID.
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// CreateRefreshToken persists a refresh token.
func (r *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

// GetValidRefreshToken returns an unexpired refresh token.
func (r *UserRepository) GetValidRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteRefreshToken removes a refresh token by its value.
func (r *UserRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token=?", token).Delete(&models.RefreshToken{}).Error
}

// DeleteRefreshTokenByID removes a refresh token by its ID.
func (r *UserRepository) DeleteRefreshTokenByID(id uint) error {
	return r.db.Delete(&models.RefreshToken{}, id).Error
}
