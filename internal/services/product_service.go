package services

import (
	"context"

	"github.com/deepakgudla/bookvault/internal/dto"
	"github.com/deepakgudla/bookvault/internal/models"
	"github.com/deepakgudla/bookvault/internal/utils"
	"gorm.io/gorm"
)

var _ ProductServiceInterface = (*ProductService)(nil)

// ProductService manages categories, products, and product images.
type ProductService struct {
	db *gorm.DB
}

// NewProductService creates a product service backed by the supplied database.
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateCategory creates a product category.
func (s *ProductService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.db.WithContext(ctx).Create(&category).Error; err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// GetCategory returns all active product categories.
func (s *ProductService) GetCategory(ctx context.Context) ([]dto.CategoryResponse, error) {
	var categories []models.Category
	if err := s.db.WithContext(ctx).Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i := range categories {
		response[i] = categoryResponse(&categories[i])
	}

	return response, nil
}

// UpdateCategory updates a product category.
func (s *ProductService) UpdateCategory(ctx context.Context, id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	var category models.Category

	if err := s.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.db.WithContext(ctx).Save(&category).Error; err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

// DeleteCategory removes a product category.
func (s *ProductService) DeleteCategory(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&models.Category{}, id).Error
}

// CreateProduct creates a product.
func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}

	if err := s.db.WithContext(ctx).Create(&product).Error; err != nil {
		return nil, err
	}

	return s.GetProduct(ctx, product.ID)
}

// GetProducts returns active products and pagination metadata.
func (s *ProductService) GetProducts(ctx context.Context, page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	var products []models.Product
	var total int64

	s.db.WithContext(ctx).Model(&models.Product{}).Where("is_active=?", true).Count(&total)

	if err := s.db.WithContext(ctx).Preload("Category").Preload("Images").
		Where("is_active=?", true).
		Offset(offset).Limit(limit).
		Find(&products).Error; err != nil {
		return nil, nil, err
	}

	response := make([]dto.ProductResponse, len(products))
	for i := range products {
		response[i] = s.convertToProductResponse(&products[i])
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return response, meta, nil
}

// GetProduct returns a product by ID.
func (s *ProductService) GetProduct(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.WithContext(ctx).Preload("Category").Preload("Images").First(&product, id).Error; err != nil {
		return nil, err
	}

	response := s.convertToProductResponse(&product)
	return &response, nil
}

// UpdateProduct updates a product.
func (s *ProductService) UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, err
	}

	return s.GetProduct(ctx, id)
}

// DeleteProduct removes a product.
func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

// AddProductImage associates an image URL with a product.
func (s *ProductService) AddProductImage(ctx context.Context, productID uint, url, altText string) error {
	var count int64
	s.db.WithContext(ctx).Model(&models.ProductImage{}).Where("product_id=?", productID).Count(&count)

	image := models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   altText,
		IsPrimary: count == 0,
	}

	return s.db.WithContext(ctx).Create(&image).Error

}

// SearchProducts uses full text search to search prodiucts
func (s *ProductService) SearchProducts(ctx context.Context, req *dto.SearchProductRequest) ([]dto.ProductSearchResult, *utils.PaginationMeta, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	query := s.db.WithContext(ctx).Model(&models.Product{}).
		Select("products.*, ts_rank(search_vector, plainto_tsquery('english', ?)) as rank", req.Query).
		Where("search_vector @@ plainto_tsquery('english', ?)", req.Query).
		Where("is_active=?", true)

	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	if req.MinPrice != nil {
		query = query.Where("price >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		query = query.Where("price <= ?", *req.MaxPrice)
	}

	var total int64
	query.Count(&total)

	type ProductsWithRank struct {
		models.Product
		Rank float32 `gorm:"column:rank"`
	}

	var rows []ProductsWithRank
	if err := query.
		Order("rank DESC, created_at DESC").
		Preload("Category").
		Preload("Images").
		Offset(offset).
		Limit(req.Limit).
		Find(&rows).Error; err != nil {
		return nil, nil, err
	}

	results := make([]dto.ProductSearchResult, len(rows))
	for i := range rows {
		results[i] = dto.ProductSearchResult{
			ProductResponse: s.convertToProductResponse(&rows[i].Product),
			Rank:            rows[i].Rank,
		}
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	meta := &utils.PaginationMeta{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return results, meta, nil
}

func (s *ProductService) convertToProductResponse(product *models.Product) dto.ProductResponse {
	return productResponse(product)
}
