package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"project-2026-06-misoastory-be-go/internal/database"
	"project-2026-06-misoastory-be-go/internal/utils"
)

// RequireAuth validates the JWT token in the Authorization header or Cookie
func RequireAuth() gin.HandlerFunc {
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			return
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Store user and position ID in context for downstream handlers and middlewares
		c.Set("user_id", claims.UserID)
		c.Set("position_id", claims.PositionID)

		c.Next()
	}
}

// RequirePermission checks if the authenticated user's position has the exact permission
func RequirePermission(resource string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		positionID, exists := c.Get("position_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
			return
		}

		var count int64
		// Join PositionPermission with Permission to check if the specific resource/action is granted
		err := database.DB.Table("position_permissions").
			Joins("JOIN permissions ON permissions.id = position_permissions.permission_id").
			Where("position_permissions.position_id = ? AND permissions.resource = ? AND permissions.action = ?", positionID, resource, action).
			Count(&count).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Database error while verifying permissions"})
			return
		}

		if count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
			return
		}

		c.Next()
	}
}
