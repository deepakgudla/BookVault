package services

import "errors"

var (
	errInvalidCredentials     = errors.New("invalid credentials")
	errInvalidRefreshToken    = errors.New("invalid refresh token")
	errRegistrationNotAllowed = errors.New("you cannot register with this")
	errProductNotFound        = errors.New("product not found")
	errInsufficientStock      = errors.New("insufficient stock")
	errCartItemNotFound       = errors.New("cart item not found")
	errCartNotFound           = errors.New("cart not found")
	errEmptyCart              = errors.New("no items found in the cart")
)
