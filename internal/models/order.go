package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents a customer's purchase order.
type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Status      OrderStatus    `json:"order_status"`
	TotalAmount float64        `json:"total_amount" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User       User        `json:"user"`
	OrderItems []OrderItem `json:"order_items"`
}

// OrderStatus describes the current state of an order.
type OrderStatus string

const (
	// OrderStatusPending indicates that an order is awaiting processing.
	OrderStatusPending OrderStatus = "pending"
	// OrderStatusConfirmed indicates that an order has been confirmed.
	OrderStatusConfirmed OrderStatus = "confirmed"
	// OrderStatusShipped indicates that an order has been shipped.
	OrderStatusShipped OrderStatus = "shipped"
	// OrderStatusDelivered indicates that an order has been delivered.
	OrderStatusDelivered OrderStatus = "delivered"
	// OrderStatusCancelled indicates that an order has been cancelled.
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents a product and quantity within an order.
type OrderItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OrderID   uint           `json:"order_id" gorm:"not null"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	Price     float64        `json:"price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Order   Order   `json:"-"`
	Product Product `json:"product"`
}

// Cart represents the active shopping cart for a user.
type Cart struct {
	ID        uint           `json:"id" gorm:"PrimaryKey"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	CartItems []CartItem `json:"cart_items"`
}

// CartItem represents a product and quantity within a cart.
type CartItem struct {
	ID        uint           `json:"id" gorm:"PrimaryKey"`
	CartID    uint           `json:"cart_id" gorm:"uniqueIndex;not null"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Cart    Cart    `json:"-"`
	Product Product `json:"product"`
}
