package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-2026-06-misoastory-be-go/internal/dto"
	"project-2026-06-misoastory-be-go/internal/services"
	"project-2026-06-misoastory-be-go/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

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
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.Response[dto.AuthResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to register user", err.Error())
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
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.Response[dto.AuthResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication failed", err.Error())
		return
	}

	// Set JWT as an HttpOnly Cookie for automatic browser/swagger auth
	c.SetCookie("token", res.Token, 3600*24, "/", "", false, true)

	utils.SuccessResponse(c, http.StatusOK, "Login successful", res)
}
