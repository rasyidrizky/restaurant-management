package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"project-2026-06-misoastory-be-go/internal/dto"
	"project-2026-06-misoastory-be-go/internal/utils"
)

// HealthCheck godoc
// @Summary Check API Health
// @Description Returns the current status of the API
// @Tags health
// @Produce json
// @Success 200 {object} dto.Response[dto.HealthResponse]
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "API is running smoothly", dto.HealthResponse{
		Status:    "OK",
		Timestamp: time.Now(),
	})
}
