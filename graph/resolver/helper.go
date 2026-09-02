package resolver

import (
	"context"
	"errors"

	"github.com/deepakgudla/bookvault/internal/utils"
)

var (
	// ErrUnauthorized indicates that the request has no valid authenticated user.
	ErrUnauthorized = errors.New("unauthorized")
)

// GetUserIDFromContext returns the authenticated user's ID from the request context.
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID := ctx.Value(utils.UserIDKey)
	if userID == nil {
		return 0, ErrUnauthorized
	}

	if id, ok := userID.(uint); ok {
		return id, nil
	}

	return 0, ErrUnauthorized
}

// GetUserRoleFromContext returns the authenticated user's role from the request context.
func GetUserRoleFromContext(ctx context.Context) (string, error) {
	userRole := ctx.Value(utils.UserRoleKey)
	if userRole == nil {
		return "", ErrUnauthorized
	}

	if role, ok := userRole.(string); ok {
		return role, nil
	}

	return "", ErrUnauthorized
}

// IsAdminFromContext reports whether the request belongs to an administrator.
func IsAdminFromContext(ctx context.Context) bool {
	role, err := GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return role == "admin"

}

// GetPagingNumbers applies default and minimum values to optional pagination inputs.
func GetPagingNumbers(page, limit *int) (pageNUmber, pageLimit int) {
	p := 1
	l := 10

	if page != nil {
		p = *page
	}

	if limit != nil {
		l = *limit
	}

	if p <= 0 {
		p = 1
	}

	if l <= 0 {
		l = 10
	}

	return p, l
}
