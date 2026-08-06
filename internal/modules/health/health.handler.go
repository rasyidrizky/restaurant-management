package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	healthtypes "project-2026-06-misoastory-be-go/internal/modules/health/types"
	"project-2026-06-misoastory-be-go/internal/common/utils"
)

// HealthCheck godoc
// @Summary Check API Health
// @Description Returns the current status of the API
// @Tags health
// @Produce json
// @Success 200 {object} dto.Response[healthtypes.HealthResponse]
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "API is running smoothly", healthtypes.HealthResponse{
		Status:    "OK",
		Timestamp: time.Now(),
	})
}

// Ping godoc
// @Summary Ping the API
// @Description Returns a simple pong response
// @Tags health
// @Produce json
// @Success 200 {object} dto.Response[string]
// @Router /ping [get]
func Ping(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Pong!", "pong")
}

var _ dto.Response[any]

