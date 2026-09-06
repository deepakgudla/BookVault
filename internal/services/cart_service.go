package services

import (
	"errors"

	"github.com/deepakgudla/bookvault/internal/dto"
	"github.com/deepakgudla/bookvault/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ CartServiceInterface = (*CartService)(nil)

// CartService manages shopping carts.
type CartService struct {
	db *gorm.DB
}

// NewCartService creates a cart service backed by the supplied database.
func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

// GetCart returns a user's cart.
func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").Where("user_id=?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return s.convertToCartResponse(&cart), nil
}

// AddToCart adds a product quantity to a user's cart.
func (s *CartService) AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var product models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, req.ProductID).Error; err != nil {
			return errors.New("product not found")
		}

		if product.Stock < req.Quantity {
			return errors.New("insufficient stock")
		}

		var cart models.Cart
		if err := tx.Where("user_id=?", userID).First(&cart).Error; err != nil {
			cart = models.Cart{UserID: userID}
			if err := tx.Create(&cart).Error; err != nil {
				return err
			}
		}

		var cartItem models.CartItem
		if err := tx.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err != nil {
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
			}
			if err := tx.Create(&cartItem).Error; err != nil {
				return err
			}
		} else {
			cartItem.Quantity += req.Quantity
			if cartItem.Quantity > product.Stock {
				return errors.New("insufficient stock")
			}
			if err := tx.Save(&cartItem).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetCart(userID)
}

// UpdateCartItem changes the quantity of an item in a user's cart.
func (s *CartService) UpdateCartItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var cartItem models.CartItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Joins("JOIN carts ON cart_items.cart_id = carts.id").
			Where("cart_items.id = ? AND carts.user_id=?", itemID, userID).
			First(&cartItem).Error; err != nil {
			return errors.New("cart item not found")
		}

		var product models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, cartItem.ProductID).Error; err != nil {
			return errors.New("product not found")
		}

		if product.Stock < req.Quantity {
			return errors.New("insufficient stock")
		}

		cartItem.Quantity = req.Quantity
		return tx.Save(&cartItem).Error
	})
	if err != nil {
		return nil, err
	}

	return s.GetCart(userID)

}

// RemoveFromCart removes an item from a user's cart.
func (s *CartService) RemoveFromCart(userID, itemID uint) error {
	result := s.db.Where("id = ? AND cart_id IN (?)", itemID, s.db.Select("id").Table("carts").Where("user_id = ?", userID)).Delete(&models.CartItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cart item not found")
	}
	return nil

}

func (s *CartService) convertToCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i := range cart.CartItems {
		subtotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subtotal

		cartItems[i] = dto.CartItemResponse{
			ID: cart.CartItems[i].ID,
			Product: dto.ProductResponse{
				ID:          cart.CartItems[i].Product.ID,
				CategoryID:  cart.CartItems[i].Product.CategoryID,
				Name:        cart.CartItems[i].Product.Name,
				Description: cart.CartItems[i].Product.Description,
				Price:       cart.CartItems[i].Product.Price,
				Stock:       cart.CartItems[i].Product.Stock,
				SKU:         cart.CartItems[i].Product.SKU,
				IsActive:    cart.CartItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          cart.CartItems[i].Product.Category.ID,
					Name:        cart.CartItems[i].Product.Category.Name,
					Description: cart.CartItems[i].Product.Category.Description,
					IsActive:    cart.CartItems[i].Product.Category.IsActive,
				},
			},
			Quantity:  cart.CartItems[i].Quantity,
			SubTotal:  subtotal,
			CreatedAt: cart.CartItems[i].CreatedAt,
			UpdatedAt: cart.CartItems[i].UpdatedAt,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
