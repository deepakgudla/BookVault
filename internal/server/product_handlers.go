package server

import (
	"strconv"

	"github.com/deepakgudla/bookvault/internal/dto"
	"github.com/deepakgudla/bookvault/internal/utils"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new category
// @Description create a new product category (Admin Only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} utils.Response{data=dto.CategoryResponse} "successfully created category"
// @Failure 400 {object} utils.Response "invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "admin access required"
// @Router /categories [post]
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

// @Summary Get All Categories
// @Description retrieve all categories that are active
// @Tags Categories
// @Produce json
// @Success 200 {object} utils.Response{data=[]dto.CategoryResponse} "categories fetched successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /categories [get]
func (s *Server) getCategories(c *gin.Context) {
	categories, err := s.productService.GetCategory()
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to create categories", err)
		return
	}

	utils.SuccessResponse(c, "categories fetched successfully", categories)
}

// @Summary Update a category
// @Description Update an existing category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} utils.Response{data=dto.CategoryResponse} "successfully updated category"
// @Failure 400 {object} utils.Response "invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "admin access required"
// @Router /categories/{id} [put]
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

// @Summary Delete a category
// @Description delete a category (Admin only)
// @Tags Categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} utils.Response "successfully deleted category"
// @Failure 400 {object} utils.Response "invalid category ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "admin access required"
// @Router /categories/{id} [delete]
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

// @Summary Create a new product
// @Description create a new product (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProductRequest true "Product data"
// @Success 201 {object} utils.Response{data=dto.ProductResponse} "successfully created product"
// @Failure 400 {object} utils.Response "invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "admin access required"
// @Router /products [post]
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

// @Summary Get all products
// @Description get the list of active products
// @Tags Products
// @Produce json
// @Param page query int false "Page Number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.ProductResponse} "successfully fetched products"
// @Failure 500 {object} utils.Response "internal server error"
// @Router /products [get]
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

// @Summary Get a product by ID
// @Description fetch detailed information about a specific product
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response{data=dto.ProductResponse} "successfully fetched product"
// @Failure 400 {object} utils.Response "invalid product ID"
// @Failure 404 {object} utils.Response "product not found"
// @Router /products/{id} [get]
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

// @Summary update a product
// @Description update an existing product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "update product data"
// @Success 200 {object} utils.Response{data=dto.ProductResponse} "successfully fetched product"
// @Failure 400 {object} utils.Response "invalid request data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "admin access required"
// @Router /products/{id} [put]
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

// @Summary Delete a product
// @Description delete a product (admin only)
// @Tags Products
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} utils.Response "successfully deleted product"
// @Failure 400 {object} utils.Response "invalid product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "admin access required"
// @Router /products/{id} [delete]
func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid category ID", err)
		return
	}

	if err := s.productService.DeleteProduct(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, "successfully deleted category", nil)
}

// @Summary Upload product Image
// @Description upload an image for a product (admin only)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param image formData file true "Image file"
// @Success 200 {object} utils.Response{data=map[string]string} "successfully uploaded image"
// @Failure 400 {object} utils.Response "invalid request or file"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "admin access required"
// @Router /products/{id}/images [post]
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

// @Summary Search Products
// @Description Search Products using full text search with ranking
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param page query int false "Page Number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category_id query int false "Filter by category ID"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.ProductSearchResult} "Search results"
// @Failure 400 {object} utils.Response "Invalid Search query"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /search [get]
func (s *Server) searchProducts(c *gin.Context) {
	var req dto.SearchProductRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequestResponse(c, "invalid search parameters", err)
		return
	}

	results, meta, err := s.productService.SearchProducts(&req)
	if err != nil {
		s.logger.Error().Err(err).Msg("Product search failed   ")
		utils.InternalServerErrorResponse(c, "search failed", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "successfully completed searching", results, *meta)
}
