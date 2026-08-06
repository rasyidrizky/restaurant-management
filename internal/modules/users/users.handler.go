package users

import (
	"net/http"
	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/common/dto"
	userdto "project-2026-06-misoastory-be-go/internal/modules/users/dto"
	"project-2026-06-misoastory-be-go/internal/common/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler processes HTTP requests for Users
type UserHandler struct {
	userService *UserService
}

// NewUserHandler acts as the constructor for UserHandler
func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes defines the API endpoints for this module
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup, m *middleware.AuthMiddleware) {
	users := router.Group("/users")
	{
		users.GET("", m.RequireAuth(), m.RequirePermission("USER", "VIEW"), h.GetUsers)
	}
}

// @Success 200 {object} dto.Response[[]userdto.UserResponse]
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve users", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", userdto.MapToUserResponses(users))
}

var _ dto.Response[any]
