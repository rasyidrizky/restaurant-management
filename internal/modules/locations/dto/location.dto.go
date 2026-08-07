package dto

import (
	"time"
	
	"project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/dto"
	
	"github.com/jinzhu/copier"
)

type LocationQuery struct {
	dto.PaginationQuery
}

type CreateLocationRequest struct {
	Name                string   `json:"name" binding:"required"`
	Address             string   `json:"address" binding:"required"`
	City                string   `json:"city" binding:"required"`
	Phone               *string  `json:"phone,omitempty"`
	Latitude            *float64 `json:"latitude,omitempty"`
	Longitude           *float64 `json:"longitude,omitempty"`
	OpeningTime         *string  `json:"opening_time,omitempty"`
	ClosingTime         *string  `json:"closing_time,omitempty"`
	IsActive            bool     `json:"is_active,omitempty"`
	IsOpen24Hours       bool     `json:"is_open_24_hours,omitempty"`
	HasDineIn           bool     `json:"has_dine_in,omitempty"`
	SupportsHomeService bool     `json:"supports_home_service,omitempty"`
	SupportsEvents      bool     `json:"supports_events,omitempty"`
}

type UpdateLocationRequest struct {
	Name                *string  `json:"name,omitempty"`
	Address             *string  `json:"address,omitempty"`
	City                *string  `json:"city,omitempty"`
	Phone               *string  `json:"phone,omitempty"`
	Latitude            *float64 `json:"latitude,omitempty"`
	Longitude           *float64 `json:"longitude,omitempty"`
	OpeningTime         *string  `json:"opening_time,omitempty"`
	ClosingTime         *string  `json:"closing_time,omitempty"`
	IsActive            *bool    `json:"is_active,omitempty"`
	IsOpen24Hours       *bool    `json:"is_open_24_hours,omitempty"`
	HasDineIn           *bool    `json:"has_dine_in,omitempty"`
	SupportsHomeService *bool    `json:"supports_home_service,omitempty"`
	SupportsEvents      *bool    `json:"supports_events,omitempty"`
}

type LocationResponse struct {
	ID                  uint      `json:"id"`
	Name                string    `json:"name"`
	Slug                string    `json:"slug"`
	Address             string    `json:"address"`
	City                string    `json:"city"`
	Phone               *string   `json:"phone,omitempty"`
	Latitude            *float64  `json:"latitude,omitempty"`
	Longitude           *float64  `json:"longitude,omitempty"`
	OpeningTime         *string   `json:"opening_time,omitempty"`
	ClosingTime         *string   `json:"closing_time,omitempty"`
	IsActive            bool      `json:"is_active"`
	IsOpen24Hours       bool      `json:"is_open_24_hours"`
	HasDineIn           bool      `json:"has_dine_in"`
	SupportsHomeService bool      `json:"supports_home_service"`
	SupportsEvents      bool      `json:"supports_events"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func MapToLocationResponse(location *models.Location) LocationResponse {
	var resp LocationResponse
	copier.Copy(&resp, location)
	return resp
}

func MapToLocationResponses(locations []models.Location) []LocationResponse {
	var responses []LocationResponse
	copier.Copy(&responses, &locations)
	if responses == nil {
		return []LocationResponse{}
	}
	return responses
}
