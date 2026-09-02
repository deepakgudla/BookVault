package repository

import "github.com/deepakgudla/bookvault/internal/models"

// UserRepositoryInterface defines user persistence operations.
type UserRepositoryInterface interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmailAndActive(email string, isActive bool) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error

	CreateRefreshToken(token *models.RefreshToken) error
	GetValidRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByID(id uint) error
}

// CartRepositoryInterface defines cart persistence operations.
type CartRepositoryInterface interface {
	GetByUserID(userID uint) (*models.Cart, error)
	Create(cart *models.Cart) error
	Update(cart *models.Cart) error
	Delete(id uint) error
}
