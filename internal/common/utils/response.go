package utils

import (
	"project-2026-06-misoastory-be-go/internal/common/dto"

	"github.com/gin-gonic/gin"
)

// SuccessResponse formats and sends a successful generic response
func SuccessResponse(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, dto.Response[any]{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse formats and sends an error generic response
func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, dto.ErrorResponse{
		Code:    statusCode,
		Message: message,
		Error:   err,
	})
}

// SuccessPaginatedResponse formats and sends a paginated successful response
func SuccessPaginatedResponse(c *gin.Context, statusCode int, message string, data any, meta dto.Meta) {
	c.JSON(statusCode, dto.PaginatedResponse[any]{
		Code:    statusCode,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}
