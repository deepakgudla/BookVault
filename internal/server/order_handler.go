package server

import (
	"strconv"

	_ "github.com/deepakgudla/bookvault/internal/dto"
	"github.com/deepakgudla/bookvault/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Create an Order
// @Description create an order from the current user's cart
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.OrderResponse} "successfully created order"
// @Failure 400 {object} utils.Response "empty cart or insufficient stock"
// @Failure 401 {object} utils.Response "unauthorized"
// @Router /orders [post]
func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	order, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "failed to create order", err)
		return
	}

	utils.CreateResponse(c, "successfully created order", order)
}

// @Summary Get user's orders
// @Description get the list of orders by user
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page Number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 201 {object} utils.PaginatedResponse{data=[]dto.OrderResponse} "orders retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders [get]
func (s *Server) getOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, meta, err := s.orderService.GetOrders(userID, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to fetch orders", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "successfully retrieved orders", &orders, *meta)
}

// @Summary Get order by ID
// @Description get the information about the specific order
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response{data=dto.OrderResponse} "order retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "order not found"
// @Router /orders/{id} [get]
func (s *Server) getOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid orderID", err)
		return
	}

	order, err := s.orderService.GetOrder(userID, uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "order not found")
		return
	}

	utils.SuccessResponse(c, "successfully retrieved order", order)
}
