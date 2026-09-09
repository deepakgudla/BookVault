package repository

import (
	"context"

	"github.com/deepakgudla/bookvault/internal/models"
)

// UserRepositoryInterface defines user persistence operations.
type UserRepositoryInterface interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmailAndActive(ctx context.Context, email string, isActive bool) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error

	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetValidRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByID(ctx context.Context, id uint) error
}

// CartRepositoryInterface defines cart persistence operations.
type CartRepositoryInterface interface {
	GetByUserID(ctx context.Context, userID uint) (*models.Cart, error)
	Create(ctx context.Context, cart *models.Cart) error
	Update(ctx context.Context, cart *models.Cart) error
	Delete(ctx context.Context, id uint) error
}
