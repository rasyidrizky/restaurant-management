package careers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/common/utils"
	careerdto "project-2026-06-misoastory-be-go/internal/modules/careers/dto"
)

// CareerHandler processes HTTP requests for Careers.
type CareerHandler struct {
	careerService *CareerService
}

// NewCareerHandler acts as the constructor for CareerHandler, injecting the CareerService.
func NewCareerHandler(careerService *CareerService) *CareerHandler {
	return &CareerHandler{careerService: careerService}
}

// RegisterRoutes defines the API endpoints and maps them to their respective handler functions.
func (h *CareerHandler) RegisterRoutes(router *gin.RouterGroup, m *middleware.AuthMiddleware) {
	careers := router.Group("/careers")
	{
		// Public read
		careers.GET("", h.GetJobPosts)
		careers.GET("/slug/:slug", h.GetJobPostBySlug)

		// Protected write
		careers.POST("", m.RequireAuth(), m.RequirePermission("JOB_POST", "ADD"), h.CreateJobPost)
		careers.PATCH("/:id", m.RequireAuth(), m.RequirePermission("JOB_POST", "UPDATE"), h.UpdateJobPost)
		careers.DELETE("/:id", m.RequireAuth(), m.RequirePermission("JOB_POST", "DELETE"), h.DeleteJobPost)
	}
}

// CreateJobPost godoc
// @Summary Create a job post
// @Description Create a new job post
// @Tags Careers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body careerdto.CreateJobPostRequest true "Job Post data"
// @Success 201 {object} dto.Response[careerdto.JobPostResponse]
// @Router /careers [post]
func (h *CareerHandler) CreateJobPost(c *gin.Context) {
	var req careerdto.CreateJobPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	jobPost, err := h.careerService.CreateJobPost(&req)
	if err != nil {
		if errors.Is(err, ErrJobPostConflict) {
			c.Error(utils.NewAppError(http.StatusConflict, "Job post conflict", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to create job post", err))
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Job post created successfully", careerdto.MapToJobPostResponse(jobPost))
}

// GetJobPosts godoc
// @Summary Get all job posts
// @Description Get a paginated list of job posts
// @Tags Careers
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search by title, division, or city"
// @Param sort query string false "Sort order"
// @Success 200 {object} dto.PaginatedResponse[[]careerdto.JobPostResponse]
// @Router /careers [get]
func (h *CareerHandler) GetJobPosts(c *gin.Context) {
	var q careerdto.JobPostQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid query parameters", err))
		return
	}

	jobPosts, meta, err := h.careerService.GetJobPosts(&q)
	if err != nil {
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve job posts", err))
		return
	}

	utils.SuccessPaginatedResponse(c, http.StatusOK, "Job posts retrieved successfully", careerdto.MapToJobPostResponses(jobPosts), meta)
}

// GetJobPostBySlug godoc
// @Summary Get a job post by slug
// @Description Get a job post's details by its slug
// @Tags Careers
// @Accept json
// @Produce json
// @Param slug path string true "Job Post Slug"
// @Success 200 {object} dto.Response[careerdto.JobPostResponse]
// @Router /careers/slug/{slug} [get]
func (h *CareerHandler) GetJobPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid job post slug", nil))
		return
	}

	jobPost, err := h.careerService.GetJobPostBySlug(slug)
	if err != nil {
		if errors.Is(err, ErrJobPostNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Job post not found", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve job post", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Job post retrieved successfully", careerdto.MapToJobPostResponse(jobPost))
}

// UpdateJobPost godoc
// @Summary Update a job post by ID
// @Description Update an existing job post
// @Tags Careers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Job Post ID"
// @Param request body careerdto.UpdateJobPostRequest true "Updated data"
// @Success 200 {object} dto.Response[careerdto.JobPostResponse]
// @Router /careers/{id} [patch]
func (h *CareerHandler) UpdateJobPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid job post ID", err))
		return
	}

	var req careerdto.UpdateJobPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	jobPost, err := h.careerService.UpdateJobPost(uint(id), &req)
	if err != nil {
		if errors.Is(err, ErrJobPostNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Job post not found", err))
			return
		}
		if errors.Is(err, ErrJobPostConflict) {
			c.Error(utils.NewAppError(http.StatusConflict, "Job post conflict", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to update job post", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Job post updated successfully", careerdto.MapToJobPostResponse(jobPost))
}

// DeleteJobPost godoc
// @Summary Delete a job post by ID
// @Description Delete an existing job post
// @Tags Careers
// @Security BearerAuth
// @Param id path int true "Job Post ID"
// @Success 200 {object} dto.Response[string]
// @Router /careers/{id} [delete]
func (h *CareerHandler) DeleteJobPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(utils.NewAppError(http.StatusBadRequest, "Invalid job post ID", err))
		return
	}

	if err := h.careerService.DeleteJobPost(uint(id)); err != nil {
		if errors.Is(err, ErrJobPostNotFound) {
			c.Error(utils.NewAppError(http.StatusNotFound, "Job post not found", err))
			return
		}
		c.Error(utils.NewAppError(http.StatusInternalServerError, "Failed to delete job post", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Job post deleted successfully", nil)
}

var _ dto.Response[any]
