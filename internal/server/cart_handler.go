package server

import (
	"strconv"

	"github.com/deepakgudla/BookVault/internal/dto"
	"github.com/deepakgudla/BookVault/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(userID)
	if err != nil {
		utils.NotFoundResponse(c, "cart nt found")
		return
	}

	utils.SuccessResponse(c, "successfully fetched cart", cart)
}

func (s *Server) addToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	cart, err := s.cartService.AddToCart(userID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "failed to add item to cart", err)
		return
	}

	utils.SuccessResponse(c, "item added to cart successfully", cart)
}

func (s *Server) updateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid cart item ID", err)
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	cart, err := s.cartService.UpdateCartItem(userID, uint(id), &req)
	if err != nil {

		utils.BadRequestResponse(c, "failed to update cart item", err)
		return
	}

	utils.SuccessResponse(c, "cart item updated successfully", cart)
}

func (s *Server) removeFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "invalid cart item ID", err)
		return
	}

	if err := s.cartService.RemoveFromCart(userID, uint(id)); err != nil {
		utils.BadRequestResponse(c, "failed to remove item from the cart", err)
		return
	}

	utils.SuccessResponse(c, "successfully removed item from the cart", nil)
}
