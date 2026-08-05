package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"project-2026-06-misoastory-be-go/internal/dto"
	"project-2026-06-misoastory-be-go/internal/utils"
)

type AuthMiddleware struct {
	db *gorm.DB
}

func NewAuthMiddleware(db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{db: db}
}

// RequireAuth validates the JWT token in the Authorization header or Cookie
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 1. Try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token = parts[1]
			}
		}

		// 2. If no header, try to get token from Cookie
		if token == "" {
			if cookieToken, err := c.Cookie("token"); err == nil {
				token = cookieToken
			}
		}

		// 3. If still no token, reject request
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Authorization token is required",
			})
			return
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Error:   "Invalid or expired token",
			})
			return
		}

		// Store user and position ID in context for downstream handlers and middlewares
		c.Set("user_id", claims.UserID)
		c.Set("position_id", claims.PositionID)

		c.Next()
	}
}

// RequirePermission checks if the authenticated user's position has the exact permission
func (m *AuthMiddleware) RequirePermission(resource string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		positionID, exists := c.Get("position_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized context",
			})
			return
		}

		var count int64
		// Join PositionPermission with Permission to check if the specific resource/action is granted
		err := m.db.Table("position_permissions").
			Joins("JOIN permissions ON permissions.id = position_permissions.permission_id").
			Where("position_permissions.position_id = ? AND permissions.resource = ? AND permissions.action = ?", positionID, resource, action).
			Count(&count).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Database error while verifying permissions",
			})
			return
		}

		if count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Code:    http.StatusForbidden,
				Message: "Forbidden",
				Error:   "You do not have permission to perform this action",
			})
			return
		}

		c.Next()
	}
}
