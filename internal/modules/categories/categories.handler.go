package categories

import (
	"errors"
	"net/http"
	"strconv"

	"project-2026-06-misoastory-be-go/internal/core/middleware"
	"project-2026-06-misoastory-be-go/internal/dto"
	"project-2026-06-misoastory-be-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService *CategoryService
}

func NewCategoryHandler(categoryService *CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

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
// @Param search query string false "Search by name or description"
// @Success 200 {object} dto.Response[[]dto.CategoryResponse]
// @Failure 500 {object} dto.ErrorResponse
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	search := c.Query("search")
	categories, err := h.categoryService.GetAllCategories(search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve categories", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", dto.MapToCategoryResponses(categories))
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get a specific category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} dto.Response[dto.CategoryResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve category", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", dto.MapToCategoryResponse(category))
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} dto.Response[dto.CategoryResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		if errors.Is(err, ErrCategoryConflict) {
			utils.ErrorResponse(c, http.StatusConflict, "Category conflict", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Category created successfully", dto.MapToCategoryResponse(category))
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body dto.UpdateCategoryRequest true "Category data"
// @Success 200 {object} dto.Response[dto.CategoryResponse]
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	category, err := h.categoryService.UpdateCategory(uint(id), &req)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found", err.Error())
			return
		}
		if errors.Is(err, ErrCategoryConflict) {
			utils.ErrorResponse(c, http.StatusConflict, "Category conflict", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Category updated successfully", dto.MapToCategoryResponse(category))
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID", err.Error())
		return
	}

	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Category not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
