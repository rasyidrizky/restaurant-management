package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"project-2026-06-misoastory-be-go/internal/dto"
	_ "project-2026-06-misoastory-be-go/internal/models"
	"project-2026-06-misoastory-be-go/internal/services"
	"project-2026-06-misoastory-be-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationService *services.LocationService
}

func NewLocationHandler(locationService *services.LocationService) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

// GetLocations godoc
// @Summary Get all locations
// @Description Get a list of all locations with optional search
// @Tags locations
// @Produce json
// @Param search query string false "Search by name or city"
// @Success 200 {object} dto.Response[[]models.Location]
// @Failure 500 {object} dto.ErrorResponse
// @Router /locations [get]
func (h *LocationHandler) GetLocations(c *gin.Context) {
	search := c.Query("search")
	locations, err := h.locationService.GetAllLocations(search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve locations", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Locations retrieved successfully", locations)
}

// GetLocationByID godoc
// @Summary Get location by ID
// @Description Get a specific location by its ID
// @Tags locations
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} dto.Response[models.Location]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /locations/{id} [get]
func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid location ID", err.Error())
		return
	}

	location, err := h.locationService.GetLocationByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Location not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve location", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Location retrieved successfully", location)
}

// CreateLocation godoc
// @Summary Create a new location
// @Description Create a new location
// @Tags locations
// @Accept json
// @Produce json
// @Param location body dto.CreateLocationRequest true "Location data"
// @Success 201 {object} dto.Response[dto.LocationResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	location, err := h.locationService.CreateLocation(&req)
	if err != nil {
		if errors.Is(err, services.ErrLocationConflict) {
			utils.ErrorResponse(c, http.StatusConflict, "Location conflict", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create location", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "Location created successfully", location)
}

// UpdateLocation godoc
// @Summary Update a location
// @Description Update an existing location by ID
// @Tags locations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param location body dto.UpdateLocationRequest true "Location data"
// @Success 200 {object} dto.Response[models.Location]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /locations/{id} [patch]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid location ID", err.Error())
		return
	}

	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	location, err := h.locationService.UpdateLocation(uint(id), &req)
	if err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Location not found", err.Error())
			return
		}
		if errors.Is(err, services.ErrLocationConflict) {
			utils.ErrorResponse(c, http.StatusConflict, "Location conflict", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update location", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Location updated successfully", location)
}

// DeleteLocation godoc
// @Summary Delete a location
// @Description Delete an existing location by ID
// @Tags locations
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} dto.Response[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid location ID", err.Error())
		return
	}

	if err := h.locationService.DeleteLocation(uint(id)); err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Location not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete location", err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "Location deleted successfully", nil)
}
