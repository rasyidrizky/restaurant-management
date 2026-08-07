package products

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/common/utils"
	productdto "project-2026-06-misoastory-be-go/internal/modules/products/dto"
)

// ProductHandler processes HTTP requests for Products.
// It acts as the controller layer, parsing requests and delegating to ProductService.
type ProductHandler struct {
	productService *ProductService
}

// NewProductHandler acts as the constructor for ProductHandler, injecting the ProductService.
func NewProductHandler(productService *ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
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

// CreateProduct godoc
// @Summary Create a product
// @Description Create a new product with optional locations mapping
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body productdto.CreateProductRequest true "Product data"
// @Success 201 {object} dto.Response[productdto.ProductResponse]
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req productdto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		if errors.Is(err, ErrProductConflict) {
			c.Error(utils.NewAppError(http.StatusConflict, "Product conflict", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to create product", err))
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", productdto.MapToProductResponse(product))
}

// GetProducts godoc
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
// @Success 200 {object} dto.PaginatedResponse[[]productdto.ProductResponse]
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var q productdto.ProductQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid query parameters", err))
		return
	}

	products, meta, err := h.productService.GetProducts(&q)
	if err != nil {
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve products", err))
		return
	}

	utils.SuccessPaginatedResponse(c, http.StatusOK, "Products retrieved successfully", productdto.MapToProductResponses(products), meta)
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Get a product's details by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} dto.Response[productdto.ProductResponse]
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid product ID", err))
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Product not found", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve product", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", productdto.MapToProductResponse(product))
}

// UpdateProduct godoc
// @Summary Update a product by ID
// @Description Update an existing product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body productdto.UpdateProductRequest true "Updated data"
// @Success 200 {object} dto.Response[productdto.ProductResponse]
// @Router /products/{id} [patch]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid product ID", err))
		return
	}

	var req productdto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	product, err := h.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Product not found", err))
			return
		}
		if errors.Is(err, ErrProductConflict) {
			c.Error(utils.NewAppError(http.StatusConflict, "Product conflict", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to update product", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", productdto.MapToProductResponse(product))
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete an existing product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} dto.Response[string]
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid product ID", err))
		return
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Product not found", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to delete product", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

var _ dto.Response[any]
