package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/deepakgudla/bookvault/internal/config"
	"github.com/deepakgudla/bookvault/internal/dto"
	"github.com/deepakgudla/bookvault/internal/events"
	"github.com/deepakgudla/bookvault/internal/models"
	"github.com/deepakgudla/bookvault/internal/repository"
	"github.com/deepakgudla/bookvault/internal/utils"
)

var _ AuthServiceInterace = (*AuthService)(nil)

// AuthService handles account authentication and token lifecycle operations.
type AuthService struct {
	userRepo       repository.UserRepositoryInterface
	cartRepo       repository.CartRepositoryInterface
	config         *config.Config
	eventPublisher events.Publisher
}

// NewAuthService creates an authentication service.
func NewAuthService(cfg *config.Config, eventPublisher events.Publisher, userRepo repository.UserRepositoryInterface, cartRepo repository.CartRepositoryInterface) *AuthService {
	return &AuthService{
		config:         cfg,
		eventPublisher: eventPublisher,
		userRepo:       userRepo,
		cartRepo:       cartRepo,
	}
}

// Register creates a customer account and its initial cart.
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {

	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("you cannot register with this")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	cart := models.Cart{UserID: user.ID}
	if err := s.cartRepo.Create(&cart); err != nil {
		fmt.Println("unable to creatr cart")
		_ = err
	}

	return s.generateAuthResponse(&user)
}

// Login authenticates a user and returns a token pair.
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmailAndActive(req.Email, true)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(user)
}

// RefreshToken validates a refresh token and returns a new token pair.
func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	refreshToken, err := s.userRepo.GetValidRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.userRepo.DeleteRefreshTokenByID(refreshToken.ID); err != nil {
		log.Println(err)
		_ = err
	}

	return s.generateAuthResponse(user)
}

// Logout invalidates a refresh token.
func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(refreshToken)
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}

	if err := s.userRepo.CreateRefreshToken(&refreshTokenModel); err != nil {
		log.Println(err)
		_ = err
	}

	err = s.eventPublisher.Publish("USER_LOGGED_IN", user, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("unable to publish user login event: %w", err)
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
