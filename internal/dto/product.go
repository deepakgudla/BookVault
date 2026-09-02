package dto

import "time"

// CreateCategoryRequest contains fields for creating a category.
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateCategoryRequest contains fields for updating a category.
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// CategoryResponse contains category details returned by the API.
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest contains fields for creating a product.
type CreateProductRequest struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
}

// UpdateProductRequest contains fields for updating a product.
type UpdateProductRequest struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	IsActive    *bool   `json:"is_active"`
}

// ProductResponse contains product details returned by the API.
type ProductResponse struct {
	ID          uint                   `json:"id"`
	CategoryID  uint                   `json:"category_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       float64                `json:"price"`
	Stock       int                    `json:"stock"`
	SKU         string                 `json:"sku"`
	IsActive    bool                   `json:"is_active"`
	Category    CategoryResponse       `json:"category"`
	Images      []ProductImageResponse `json:"images"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ProductImageResponse contains product image details.
type ProductImageResponse struct {
	ID        uint      `json:"id"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}
