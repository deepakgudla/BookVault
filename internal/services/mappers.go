package services

import (
	"github.com/deepakgudla/bookvault/internal/dto"
	"github.com/deepakgudla/bookvault/internal/models"
)

func categoryResponse(category *models.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID: category.ID, Name: category.Name, Description: category.Description,
		IsActive: category.IsActive, CreatedAt: category.CreatedAt, UpdatedAt: category.UpdatedAt,
	}
}

func productResponse(product *models.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		image := &product.Images[i]
		images[i] = dto.ProductImageResponse{
			ID: image.ID, URL: image.URL, AltText: image.AltText,
			IsPrimary: image.IsPrimary, CreatedAt: image.CreatedAt,
		}
	}

	return dto.ProductResponse{
		ID: product.ID, CategoryID: product.CategoryID, Name: product.Name,
		Description: product.Description, Price: product.Price, Stock: product.Stock,
		SKU: product.SKU, IsActive: product.IsActive, Category: categoryResponse(&product.Category),
		Images: images, CreatedAt: product.CreatedAt, UpdatedAt: product.UpdatedAt,
	}
}
