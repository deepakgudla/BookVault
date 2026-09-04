package services

import (
	"mime/multipart"

	"github.com/deepakgudla/bookvault/internal/dto"
	"github.com/deepakgudla/bookvault/internal/utils"
)

// AuthServiceInterace defines authentication service operations.
type AuthServiceInterace interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(refreshToken string) error
}

// UserServiceInterface defines profile operations.
type UserServiceInterface interface {
	GetProfile(userID uint) (*dto.UserResponse, error)
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
}

// ProductServiceInterface defines category and product operations.
type ProductServiceInterface interface {
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategory() ([]dto.CategoryResponse, error)
	UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(id uint) error

	CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error)
	GetProduct(id uint) (*dto.ProductResponse, error)
	UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error

	AddProductImage(productID uint, url, alText string) error
	SearchProducts(req *dto.SearchProductRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error)
}

// CartServiceInterface defines shopping cart operations.
type CartServiceInterface interface {
	GetCart(userID uint) (*dto.CartResponse, error)
	AddToCart(userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateCartItem(userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error)
	RemoveFromCart(userID, itemID uint) error
}

// OrderServiceInterface defines order operations.
type OrderServiceInterface interface {
	CreateOrder(userID uint) (*dto.OrderResponse, error)
	GetOrders(userID uint, page, limit int) ([]dto.OrderResponse, *utils.PaginationMeta, error)
	GetOrder(userID, orderID uint) (*dto.OrderResponse, error)
}

// UploadServiceInterface defines product image upload operations.
type UploadServiceInterface interface {
	UploadProductImage(productID uint, file *multipart.FileHeader) (string, error)
}
