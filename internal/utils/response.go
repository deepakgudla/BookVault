package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

// PaginatedResponse is an API response envelope with pagination metadata.
type PaginatedResponse struct {
	Response
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta describes a paginated result set.
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// SuccessResponse writes a successful HTTP response.
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// CreateResponse writes a successful resource-creation response.
func CreateResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse writes an error response with an optional underlying error.
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	c.JSON(statusCode, response)
}

// BadRequestResponse writes a bad-request response.
func BadRequestResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusBadRequest, message, err)
}

// UnauthorizedResponse writes an authentication-required response.
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

// ForbidddenResponse writes a forbidden response.
func ForbidddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

// NotFoundResponse writes a not-found response.
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// InternalServerErrorResponse writes an internal-server-error response.
func InternalServerErrorResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusInternalServerError, message, err)
}

// PaginatedSuccessResponse writes a successful response with pagination metadata.
func PaginatedSuccessResponse(c *gin.Context, message string, data interface{}, meta PaginationMeta) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Meta: meta,
	})
}
