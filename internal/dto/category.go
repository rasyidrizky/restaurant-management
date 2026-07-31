package dto

import (
	"time"

	"project-2026-06-misoastory-be-go/internal/models"
)

type CreateCategoryRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

type CategoryResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func MapToCategoryResponse(category *models.Category) CategoryResponse {
	return CategoryResponse{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		Description:  category.Description,
		DisplayOrder: category.DisplayOrder,
		IsActive:     category.IsActive,
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
	}
}

func MapToCategoryResponses(categories []models.Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = MapToCategoryResponse(&cat)
	}
	return responses
}
