package dto

import (
	"time"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	"project-2026-06-misoastory-be-go/internal/common/models"
)

// CreateJobPostRequest is the payload for creating a new job post
type CreateJobPostRequest struct {
	Title           string `json:"title" binding:"required"`
	Division        string `json:"division" binding:"required"`
	City            string `json:"city" binding:"required"`
	ExperienceLevel string `json:"experience_level" binding:"required"`
	EducationLevel  string `json:"education_level" binding:"required"`
	EmploymentType  string `json:"employment_type" binding:"required"`
	Description     string `json:"description" binding:"required"`
	Requirements    string `json:"requirements" binding:"required"`
	IsActive        *bool  `json:"is_active"`
}

// UpdateJobPostRequest is the payload for updating an existing job post
type UpdateJobPostRequest struct {
	Title           *string `json:"title"`
	Division        *string `json:"division"`
	City            *string `json:"city"`
	ExperienceLevel *string `json:"experience_level"`
	EducationLevel  *string `json:"education_level"`
	EmploymentType  *string `json:"employment_type"`
	Description     *string `json:"description"`
	Requirements    *string `json:"requirements"`
	IsActive        *bool   `json:"is_active"`
}

// JobPostQuery extends the base pagination query for searching job posts
type JobPostQuery struct {
	dto.PaginationQuery
}

// JobPostResponse is the payload returned for a job post
type JobPostResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Division        string    `json:"division"`
	City            string    `json:"city"`
	ExperienceLevel string    `json:"experience_level"`
	EducationLevel  string    `json:"education_level"`
	EmploymentType  string    `json:"employment_type"`
	Description     string    `json:"description"`
	Requirements    string    `json:"requirements"`
	IsActive        *bool     `json:"is_active"`
	PublishedAt     time.Time `json:"published_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MapToJobPostResponse maps a JobPost model to a JobPostResponse DTO
func MapToJobPostResponse(jobPost *models.JobPost) JobPostResponse {
	return JobPostResponse{
		ID:              jobPost.ID,
		Title:           jobPost.Title,
		Slug:            jobPost.Slug,
		Division:        jobPost.Division,
		City:            jobPost.City,
		ExperienceLevel: jobPost.ExperienceLevel,
		EducationLevel:  jobPost.EducationLevel,
		EmploymentType:  jobPost.EmploymentType,
		Description:     jobPost.Description,
		Requirements:    jobPost.Requirements,
		IsActive:        jobPost.IsActive,
		PublishedAt:     jobPost.PublishedAt,
		CreatedAt:       jobPost.CreatedAt,
		UpdatedAt:       jobPost.UpdatedAt,
	}
}

// MapToJobPostResponses maps a slice of JobPost models to a slice of JobPostResponse DTOs
func MapToJobPostResponses(jobPosts []models.JobPost) []JobPostResponse {
	responses := make([]JobPostResponse, len(jobPosts))
	for i, jobPost := range jobPosts {
		responses[i] = MapToJobPostResponse(&jobPost)
	}
	return responses
}
