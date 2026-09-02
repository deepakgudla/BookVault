package utils

// ContextKey is the type used for values stored in request contexts.
type ContextKey string

const (
	// UserIDKey stores the authenticated user's ID.
	UserIDKey ContextKey = "user_id"
	// UserEmailKey stores the authenticated user's email address.
	UserEmailKey ContextKey = "user_email"
	// UserRoleKey stores the authenticated user's role.
	UserRoleKey ContextKey = "user_role"
	// GinContextKey stores the Gin request context.
	GinContextKey ContextKey = "gin_context"
)
