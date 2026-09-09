package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID returns the user with the given ID.
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmailAndActive returns a user matching an email address and active state.
func (r *UserRepository) GetByEmailAndActive(ctx context.Context, email string, isActive bool) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, isActive).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create persists a new user.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update persists changes to an existing user.
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete removes a user by ID.
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// CreateRefreshToken persists a refresh token.
func (r *UserRepository) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	token.Token = hashToken(token.Token)
	return r.db.WithContext(ctx).Create(token).Error
}

// GetValidRefreshToken returns an unexpired refresh token.
func (r *UserRepository) GetValidRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ? AND expires_at > ?", hashToken(token), time.Now()).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteRefreshToken removes a refresh token by its value.
func (r *UserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token=?", hashToken(token)).Delete(&models.RefreshToken{}).Error
}

// DeleteRefreshTokenByID removes a refresh token by its ID.
func (r *UserRepository) DeleteRefreshTokenByID(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.RefreshToken{}, id).Error
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
