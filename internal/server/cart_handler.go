package server

import (
	"strconv"

	"github.com/deepakgudla/BookVault/internal/dto"
	"github.com/deepakgudla/BookVault/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user's cart
// @Description Retrieve current user's shopping cart with all items
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dto.CartResponse} "successfully retrieved cart"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "cart not found"
// @Router /cart [get]
func (s *Server) getCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := s.cartService.GetCart(userID)
	if err != nil {
		utils.NotFoundResponse(c, "cart nt found")
		return
	}

	utils.SuccessResponse(c, "successfully fetched cart", cart)
}

// @Summary Add items to cart
// @Description users can add products to their cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "item added to cart"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "successfully added products to the cart"
// @Failure 400 {object} utils.Response "invalid request data or insufficient stock"
// @Failure 401 {object} utils.Response "unauthorized"
// @Router /cart/items [post]
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

// @Summary Update items to cart
// @Description update the quantity of an item in the user's cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Param request body dto.UpdateCartItemRequest true "New quantity"
// @Success 200 {object} utils.Response{data=dto.CartResponse} "cart item updated successfully"
// @Failure 400 {object} utils.Response "invalid request data or insufficient stock"
// @Failure 401 {object} utils.Response "unauthorized"
// @Router /cart/items/{id} [put]
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

// @Summary Delete items from cart
// @Description users can delete cart
// @Tags Cart
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cart Item ID"
// @Success 200 {object} utils.Response "successfully deleted item from the cart"
// @Failure 400 {object} utils.Response "invalid cart item"
// @Failure 401 {object} utils.Response "unauthorized"
// @Router /cart/item/{id} [delete]
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
