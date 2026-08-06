package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	authdto "project-2026-06-misoastory-be-go/internal/modules/auth/dto"

	"project-2026-06-misoastory-be-go/internal/common/utils"
)

// AuthHandler processes HTTP requests for Authentication
type AuthHandler struct {
	authService *AuthService
}

// NewAuthHandler acts as the constructor for AuthHandler
func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes defines the API endpoints for this module
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authdto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.Response[authdto.AuthResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req authdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	res, err := h.authService.Register(&req)
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Failed to register user", err))
		return
	}

	// Set JWT as an HttpOnly Cookie for automatic browser/swagger auth
	c.SetCookie("token", res.Token, 3600*24, "/", "", false, true)

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", res)
}

// Login godoc
// @Summary Login
// @Description Authenticate user and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authdto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.Response[authdto.AuthResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	res, err := h.authService.Login(&req)
	if err != nil {
		c.Error(utils.NewAppError(http.StatusUnauthorized, "Authentication failed", err))
		return
	}

	// Set JWT as an HttpOnly Cookie for automatic browser/swagger auth
	c.SetCookie("token", res.Token, 3600*24, "/", "", false, true)

	utils.SuccessResponse(c, http.StatusOK, "Login successful", res)
}

var _ dto.Response[any]

