package services

import (
	"context"
	"mime/multipart"

	"github.com/deepakgudla/bookvault/internal/dto"
	"github.com/deepakgudla/bookvault/internal/utils"
)

// AuthServiceInterace defines authentication service operations.
type AuthServiceInterace interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

// UserServiceInterface defines profile operations.
type UserServiceInterface interface {
	GetProfile(ctx context.Context, userID uint) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
}

// ProductServiceInterface defines category and product operations.
type ProductServiceInterface interface {
	CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategory(ctx context.Context) ([]dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id uint) error

	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProducts(ctx context.Context, page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error)
	GetProduct(ctx context.Context, id uint) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uint) error

	AddProductImage(ctx context.Context, productID uint, url, alText string) error
	SearchProducts(ctx context.Context, req *dto.SearchProductRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error)
}

// CartServiceInterface defines shopping cart operations.
type CartServiceInterface interface {
	GetCart(ctx context.Context, userID uint) (*dto.CartResponse, error)
	AddToCart(ctx context.Context, userID uint, req *dto.AddToCartRequest) (*dto.CartResponse, error)
	UpdateCartItem(ctx context.Context, userID, itemID uint, req *dto.UpdateCartItemRequest) (*dto.CartResponse, error)
	RemoveFromCart(ctx context.Context, userID, itemID uint) error
}

// OrderServiceInterface defines order operations.
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, userID uint) (*dto.OrderResponse, error)
	GetOrders(ctx context.Context, userID uint, page, limit int) ([]dto.OrderResponse, *utils.PaginationMeta, error)
	GetOrder(ctx context.Context, userID, orderID uint) (*dto.OrderResponse, error)
}

// UploadServiceInterface defines product image upload operations.
type UploadServiceInterface interface {
	UploadProductImage(productID uint, file *multipart.FileHeader) (string, error)
}
