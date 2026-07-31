package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"project-2026-06-misoastory-be-go/internal/dto"
	"project-2026-06-misoastory-be-go/internal/services"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationService *services.LocationService
}

func NewLocationHandler() *LocationHandler {
	return &LocationHandler{
		locationService: services.NewLocationService(),
	}
}

// GetLocations godoc
// @Summary Get all locations
// @Description Get a list of all locations with optional search
// @Tags locations
// @Produce json
// @Param search query string false "Search by name or city"
// @Success 200 {object} map[string][]models.Location
// @Router /locations [get]
func (h *LocationHandler) GetLocations(c *gin.Context) {
	search := c.Query("search")
	locations, err := h.locationService.GetAllLocations(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": locations})
}

// GetLocationByID godoc
// @Summary Get location by ID
// @Description Get a specific location by its ID
// @Tags locations
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} models.Location
// @Failure 404 {object} map[string]string
// @Router /locations/{id} [get]
func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	location, err := h.locationService.GetLocationByID(uint(id))
	if err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, location)
}

// CreateLocation godoc
// @Summary Create a new location
// @Description Create a new location
// @Tags locations
// @Accept json
// @Produce json
// @Param location body dto.CreateLocationRequest true "Location data"
// @Success 201 {object} models.Location
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := h.locationService.CreateLocation(&req)
	if err != nil {
		if errors.Is(err, services.ErrLocationConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, location)
}

// UpdateLocation godoc
// @Summary Update a location
// @Description Update an existing location by ID
// @Tags locations
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param location body dto.UpdateLocationRequest true "Location data"
// @Success 200 {object} models.Location
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /locations/{id} [patch]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := h.locationService.UpdateLocation(uint(id), &req)
	if err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrLocationConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, location)
}

// DeleteLocation godoc
// @Summary Delete a location
// @Description Delete an existing location by ID
// @Tags locations
// @Produce json
// @Param id path int true "Location ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location ID"})
		return
	}

	if err := h.locationService.DeleteLocation(uint(id)); err != nil {
		if errors.Is(err, services.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.Status(http.StatusNoContent)
}
