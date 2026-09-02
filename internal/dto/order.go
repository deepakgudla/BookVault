package dto

import "time"

// AddToCartRequest identifies a product and quantity to add to a cart.
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest contains a replacement cart-item quantity.
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// CartResponse contains a cart and its calculated total.
type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	CartItems []CartItemResponse `json:"cart_items"`
	Total     float64            `json:"total"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// CartItemResponse contains a cart item and its product details.
type CartItemResponse struct {
	ID        uint            `json:"id"`
	Product   ProductResponse `json:"product"`
	Quantity  int             `json:"quantity"`
	SubTotal  float64         `json:"subtotal"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// OrderResponse contains an order and its items.
type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// OrderItemResponse contains an ordered product and quantity.
type OrderItemResponse struct {
	ID        uint            `json:"id"`
	Product   ProductResponse `json:"product"`
	Quantity  int             `json:"quantity"`
	Price     float64         `json:"price"`
	CreatedAt time.Time       `json:"created_at"`
}
