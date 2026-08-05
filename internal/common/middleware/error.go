package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"project-2026-06-misoastory-be-go/internal/common/utils"
)

// ErrorHandler is a global middleware that catches and formats all errors passed to c.Error()
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Proceed to the actual handler
		c.Next()

		// Check if there are any errors after the handler completes
		if len(c.Errors) > 0 {
			// Extract the last error
			err := c.Errors.Last().Err

			// Check if the error is our custom AppError
			var appErr *utils.AppError
			if errors.As(err, &appErr) {
				var errStr string
				if appErr.Err != nil {
					errStr = appErr.Err.Error()
				}
				utils.ErrorResponse(c, appErr.StatusCode, appErr.Message, errStr)
				return
			}

			// If it's a standard error that isn't wrapped in AppError, default to 500
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", err.Error())
		}
	}
}
