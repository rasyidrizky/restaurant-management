package locations

import (
	"errors"
	"net/http"
	"strconv"

	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/common/dto"
	locationdto "project-2026-06-misoastory-be-go/internal/modules/locations/dto"
	_ "project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/utils"

	"github.com/gin-gonic/gin"
)

// LocationHandler processes HTTP requests for Locations
type LocationHandler struct {
	locationService *LocationService
}

// NewLocationHandler acts as the constructor for LocationHandler
func NewLocationHandler(locationService *LocationService) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

// RegisterRoutes defines the API endpoints for this module
func (h *LocationHandler) RegisterRoutes(router *gin.RouterGroup, m *middleware.AuthMiddleware) {
	locations := router.Group("/locations")
	{
		// Public read
		locations.GET("", h.GetLocations)
		locations.GET("/:id", h.GetLocationByID)
		// Protected write
		locations.POST("", m.RequireAuth(), m.RequirePermission("LOCATION", "ADD"), h.CreateLocation)
		locations.PATCH("/:id", m.RequireAuth(), m.RequirePermission("LOCATION", "UPDATE"), h.UpdateLocation)
		locations.DELETE("/:id", m.RequireAuth(), m.RequirePermission("LOCATION", "DELETE"), h.DeleteLocation)
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
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve locations", err))
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
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid location ID", err))
		return
	}

	location, err := h.locationService.GetLocationByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrLocationNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Location not found", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve location", err))
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
// @Param location body locationdto.CreateLocationRequest true "Location data"
// @Success 201 {object} dto.Response[locationdto.LocationResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var req locationdto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	location, err := h.locationService.CreateLocation(&req)
	if err != nil {
		if errors.Is(err, ErrLocationConflict) {
			c.Error(utils.NewAppError(http.StatusConflict, "Location conflict", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to create location", err))
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
// @Param location body locationdto.UpdateLocationRequest true "Location data"
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
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid location ID", err))
		return
	}

	var req locationdto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	location, err := h.locationService.UpdateLocation(uint(id), &req)
	if err != nil {
		if errors.Is(err, ErrLocationNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Location not found", err))
			return
		}
		if errors.Is(err, ErrLocationConflict) {
			c.Error(utils.NewAppError(http.StatusConflict, "Location conflict", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to update location", err))
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
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid location ID", err))
		return
	}

	if err := h.locationService.DeleteLocation(uint(id)); err != nil {
		if errors.Is(err, ErrLocationNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Location not found", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to delete location", err))
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "Location deleted successfully", nil)
}

var _ dto.Response[any]

