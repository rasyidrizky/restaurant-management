package types

import (
	"project-2026-06-misoastory-be-go/internal/common/dto"
	"project-2026-06-misoastory-be-go/internal/common/models"
)

// CreateProductRequest represents the body for creating a product
type CreateProductRequest struct {
	Name         string   `json:"name" binding:"required"`
	Description  *string  `json:"description"`
	Price        float64  `json:"price" binding:"required,min=0"`
	ImageURL     *string  `json:"imageUrl"`
	CategoryID   uint     `json:"categoryId" binding:"required"`
	IsAvailable  *bool    `json:"isAvailable"`
	IsBestSeller *bool    `json:"isBestSeller"`
	DisplayOrder *int     `json:"displayOrder"`
	LocationIDs  []uint   `json:"locationIds"`
}

// UpdateProductRequest represents the body for updating a product
type UpdateProductRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Price        *float64 `json:"price" binding:"omitempty,min=0"`
	ImageURL     *string  `json:"imageUrl"`
	CategoryID   *uint    `json:"categoryId"`
	IsAvailable  *bool    `json:"isAvailable"`
	IsBestSeller *bool    `json:"isBestSeller"`
	DisplayOrder *int     `json:"displayOrder"`
	LocationIDs  []uint   `json:"locationIds"`
}

// ProductQuery extends the base pagination query
type ProductQuery struct {
	dto.PaginationQuery
	CategoryID *uint `form:"categoryId"`
	LocationID *uint `form:"locationId"`
}

// ProductResponse is the API output model
type ProductResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  *string `json:"description"`
	Price        float64 `json:"price"`
	ImageURL     *string `json:"imageUrl"`
	IsAvailable  bool    `json:"isAvailable"`
	IsBestSeller bool    `json:"isBestSeller"`
	DisplayOrder int     `json:"displayOrder"`
	CategoryID   uint    `json:"categoryId"`

	Category  *models.Category         `json:"category,omitempty"`
	Locations []models.ProductLocation `json:"locations,omitempty"`
}

// MapToProductResponse maps a model to a response
func MapToProductResponse(product *models.Product) *ProductResponse {
	if product == nil {
		return nil
	}
	return &ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Slug:         product.Slug,
		Description:  product.Description,
		Price:        product.Price,
		ImageURL:     product.ImageURL,
		IsAvailable:  product.IsAvailable,
		IsBestSeller: product.IsBestSeller,
		DisplayOrder: product.DisplayOrder,
		CategoryID:   product.CategoryID,
		Category:     product.Category,
		Locations:    product.Locations,
	}
}

// MapToProductResponses maps a slice of models to responses
func MapToProductResponses(products []models.Product) []ProductResponse {
	responses := make([]ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *MapToProductResponse(&product)
	}
	return responses
}
