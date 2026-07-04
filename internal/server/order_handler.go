package server

import (
	"strconv"

	"github.com/deepakgudla/BookVault/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	order, err := s.orderService.CreateOrder(int(userID))
	if err != nil {
		utils.BadRequestResponse(c, "failed to create order", err)
		return
	}

	utils.CreateResponse(c, "successfully created order", order)
}

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
