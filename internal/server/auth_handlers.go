package server

import (
	"github.com/deepakgudla/BookVault/internal/dto"
	"github.com/deepakgudla/BookVault/internal/utils"

	"github.com/gin-gonic/gin"
)

// @Summary Register a new user
func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	response, err := s.authService.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "registration failed", err)
		return
	}

	utils.CreateResponse(c, "user registered successfully", response)
}

func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}
	response, err := s.authService.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "login failed")
		return
	}

	utils.SuccessResponse(c, "login successful", response)
}

func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	response, err := s.authService.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "refresh Token failed")
		return
	}

	utils.SuccessResponse(c, "refresh Token successful", response)
}

func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	if err := s.authService.Logout(req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, "logout failed", err)
		return
	}

	utils.SuccessResponse(c, "login successful", nil)
}

func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	profile, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "user not found")
		return
	}

	utils.SuccessResponse(c, "user profile fetched successfully", profile)
}

func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	profile, err := s.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "failed to update profile", err)
		return
	}

	utils.SuccessResponse(c, "profile updated successfully", profile)
}
