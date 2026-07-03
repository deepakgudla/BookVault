package server

import (
	"strconv"

	"github.com/deepakgudla/BookVault/internal/dto"
	"github.com/deepakgudla/BookVault/internal/utils"

	"github.com/gin-gonic/gin"
)

func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	category, err := s.productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to create category", err)
		return
	}

	utils.CreateResponse(c, "successfully created category", category)
}

func (s *Server) getCategories(c *gin.Context) {
	categories, err := s.productService.GetCategory()
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to create categories", err)
		return
	}

	utils.SuccessResponse(c, "categories fetched successfully", categories)
}

func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid category ID", err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	category, err := s.productService.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to update category", err)
		return
	}

	utils.SuccessResponse(c, "successfully updated category", category)
}

func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid category ID", err)
		return
	}

	if err := s.productService.DeleteCategory(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "failed to delete category", err)
		return
	}

	utils.SuccessResponse(c, "successfully deleted category", nil)
}

func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	product, err := s.productService.CreateProduct(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to create product", err)
		return
	}

	utils.SuccessResponse(c, "successfully created product", product)
}

func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, meta, err := s.productService.GetProducts(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to fetch products", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "successfully fetched products", products, *meta)
}

func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid product ID", err)
		return
	}

	product, err := s.productService.GetProduct(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "product not found")
		return
	}

	utils.SuccessResponse(c, "product fetched successfully", product)
}

func (s *Server) updateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid product ID", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	product, err := s.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to update product", err)
		return
	}

	utils.SuccessResponse(c, "successfully updated product", product)
}

func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid category ID", err)
	}

	if err := s.productService.DeleteProduct(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "successfully deleted category", nil)
}

func (s *Server) uploadProductImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid product ID", err)
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequestResponse(c, "no file has been uploaded", err)
		return
	}

	// uploadProvider := providers.NewLocalUploadProvider(s.config.Upload.Path)
	// uploadService := services.NewUploadService(uploadProvider)

	url, err := s.uploadService.UploadProductImage(uint(id), file)
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to upload Image", err)
		return
	}

	if err := s.productService.AddProductImage(uint(id), url, file.Filename); err != nil {
		utils.InternalServerErrorResponse(c, "failed to save image record", err)
		return
	}

	utils.SuccessResponse(c, "image has been uploaded successfully", map[string]string{"url": url})

}
