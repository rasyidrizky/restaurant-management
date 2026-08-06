package categories

import (
	"errors"
	"net/http"
	"strconv"

	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/common/dto"
	categorytypes "project-2026-06-misoastory-be-go/internal/modules/categories/types"
	"project-2026-06-misoastory-be-go/internal/common/utils"

	"github.com/gin-gonic/gin"
)

// CategoryHandler processes HTTP requests for Categories
type CategoryHandler struct {
	categoryService *CategoryService
}

// NewCategoryHandler acts as the constructor for CategoryHandler
func NewCategoryHandler(categoryService *CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// RegisterRoutes defines the API endpoints for this module
func (h *CategoryHandler) RegisterRoutes(router *gin.RouterGroup, m *middleware.AuthMiddleware) {
	categories := router.Group("/categories")
	{
		// Public read
		categories.GET("", h.GetCategories)
		categories.GET("/:id", h.GetCategoryByID)
		// Protected write
		categories.POST("", m.RequireAuth(), m.RequirePermission("CATEGORY", "ADD"), h.CreateCategory)
		categories.PATCH("/:id", m.RequireAuth(), m.RequirePermission("CATEGORY", "UPDATE"), h.UpdateCategory)
		categories.DELETE("/:id", m.RequireAuth(), m.RequirePermission("CATEGORY", "DELETE"), h.DeleteCategory)
	}
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get a list of all categories with optional search
// @Tags categories
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param search query string false "Search term"
// @Param sort query string false "Sort order"
// @Success 200 {object} dto.PaginatedResponse[[]categorytypes.CategoryResponse]
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var req categorytypes.GetAllCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid query parameters", err))
		return
	}

	categories, meta, err := h.categoryService.GetAllCategories(&req)
	if err != nil {
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve categories", err))
		return
	}

	utils.SuccessPaginatedResponse(c, http.StatusOK, "Categories retrieved successfully", categorytypes.MapToCategoryResponses(categories), meta)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a specific category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.Response[categorytypes.CategoryResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid category ID", err))
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Category not found", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve category", err))
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", categorytypes.MapToCategoryResponse(category))
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body categorytypes.CreateCategoryRequest true "Category data"
// @Success 201 {object} dto.Response[categorytypes.CategoryResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req categorytypes.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		if errors.Is(err, ErrCategoryConflict) {
			c.Error(utils.NewAppError(http.StatusConflict, "Category conflict", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to create category", err))
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", categorytypes.MapToCategoryResponse(category))
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body categorytypes.UpdateCategoryRequest true "Category data"
// @Success 200 {object} dto.Response[categorytypes.CategoryResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /categories/{id} [patch]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid category ID", err))
		return
	}

	var req categorytypes.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), &req)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Category not found", err))
			return
		}
		if errors.Is(err, ErrCategoryConflict) {
			c.Error(utils.NewAppError(http.StatusConflict, "Category conflict", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to update category", err))
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", categorytypes.MapToCategoryResponse(category))
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete an existing category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.Response[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid category ID", err))
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Category not found", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to delete category", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}

var _ dto.Response[any]

