package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents an account in the system.
type User struct {
	ID        uint           `json:"id" gorm:"primary key"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Phone     string         `json:"phone"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	Role      UserRole       `json:"role" gorm:"default:customer"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	RefreshTokens []RefreshToken `json:"-"`
	Orders        []Order        `json:"-"`
	Cart          Cart           `json:"-"`
}

// UserRole identifies the permissions assigned to a user.
type UserRole string

const (
	// UserRoleCustomer identifies a standard customer account.
	UserRoleCustomer UserRole = "customer"
	// UserRoleAdmin identifies an administrator account.
	UserRoleAdmin UserRole = "admin"
)

// RefreshToken represents a token used to refresh an authenticated session.
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primary key"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Token     string         `json:"token" gorm:"UniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"-"`
}
