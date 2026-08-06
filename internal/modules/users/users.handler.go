package users

import (
	"net/http"
	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/common/dto"
	usertypes "project-2026-06-misoastory-be-go/internal/modules/users/types"
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

// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "Search term"
// @Param sort query string false "Sort order"
// @Success 200 {object} dto.PaginatedResponse[[]usertypes.UserResponse]
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [get]
// @Security BearerAuth
func (h *UserHandler) GetUsers(c *gin.Context) {
	var req usertypes.GetAllUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid query parameters", err))
		return
	}

	users, meta, err := h.userService.GetUsers(&req)
	if err != nil {
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve users", err))
		return
	}

	utils.SuccessPaginatedResponse(c, http.StatusOK, "Users retrieved successfully", usertypes.MapToUserResponses(users), meta)
}

var _ dto.Response[any]
