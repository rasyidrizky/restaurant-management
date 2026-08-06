package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/common/utils"
	productstypes "project-2026-06-misoastory-be-go/internal/modules/products/types"
)

// ProductHandler processes HTTP requests for Products.
type ProductHandler struct {
	service *ProductService
}

// NewProductHandler acts as the constructor for ProductHandler, injecting the ProductService.
func NewProductHandler(service *ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// RegisterRoutes defines the API endpoints and maps them to their respective handler functions.
func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup, m *middleware.AuthMiddleware) {
	products := router.Group("/products")
	{
		// Public read
		products.GET("", h.GetProducts)
		products.GET("/:id", h.GetProductByID)

		// Protected write
		products.POST("", m.RequireAuth(), m.RequirePermission("PRODUCT", "ADD"), h.CreateProduct)
		products.PATCH("/:id", m.RequireAuth(), m.RequirePermission("PRODUCT", "UPDATE"), h.UpdateProduct)
		products.DELETE("/:id", m.RequireAuth(), m.RequirePermission("PRODUCT", "DELETE"), h.DeleteProduct)
	}
}

// Create godoc
// @Summary Create a product
// @Description Create a new product with optional locations mapping
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body productstypes.CreateProductRequest true "Product data"
// @Success 201 {object} dto.Response[productstypes.ProductResponse]
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req productstypes.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	product, err := h.service.CreateProduct(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", productstypes.MapToProductResponse(product))
}

// FindAll godoc
// @Summary Get all products
// @Description Get a paginated list of products
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by name"
// @Param sort query string false "Sort order"
// @Param categoryId query int false "Filter by category ID"
// @Param locationId query int false "Filter by location ID"
// @Success 200 {object} dto.PaginatedResponse[[]productstypes.ProductResponse]
// @Router /api/v1/products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var q productstypes.ProductQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	products, meta, err := h.service.GetProducts(&q)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve products", err.Error())
		return
	}

	utils.SuccessPaginatedResponse(c, http.StatusOK, "Products retrieved successfully", productstypes.MapToProductResponses(products), meta)
}

// FindOne godoc
// @Summary Get a product by ID
// @Description Get a product's details by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.Response[productstypes.ProductResponse]
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Product not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", productstypes.MapToProductResponse(product))
}

// Update godoc
// @Summary Update a product by ID
// @Description Update an existing product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body productstypes.UpdateProductRequest true "Updated data"
// @Success 200 {object} dto.Response[productstypes.ProductResponse]
// @Router /api/v1/products/{id} [patch]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	var req productstypes.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	product, err := h.service.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", productstypes.MapToProductResponse(product))
}

// Delete godoc
// @Summary Delete a product by ID
// @Description Delete an existing product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} dto.Response[string]
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

var _ dto.Response[any]
