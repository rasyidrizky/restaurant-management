package categorytypes

import (
	"time"

	"project-2026-06-misoastory-be-go/internal/common/models"
	
	"github.com/jinzhu/copier"
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
	var resp CategoryResponse
	copier.Copy(&resp, category)
	return resp
}

func MapToCategoryResponses(categories []models.Category) []CategoryResponse {
	var responses []CategoryResponse
	copier.Copy(&responses, &categories)
	if responses == nil {
		return []CategoryResponse{}
	}
	return responses
}
