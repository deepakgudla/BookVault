package repository

import (
	"github.com/deepakgudla/bookvault/internal/models"
	"gorm.io/gorm"
)

// CartRepository persists shopping carts.
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a cart repository backed by the supplied database.
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

// GetByUserID returns the cart belonging to a user.
func (r *CartRepository) GetByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// Create persists a new cart.
func (r *CartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

// Update persists changes to an existing cart.
func (r *CartRepository) Update(cart *models.Cart) error {
	return r.db.Save(cart).Error
}

// Delete removes a cart by ID.
func (r *CartRepository) Delete(id uint) error {
	return r.db.Delete(&models.Cart{}, id).Error
}
