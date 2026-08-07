package models

import (
	"time"
)

// JobPost represents a career opportunity posted by the restaurant.
type JobPost struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Title           string    `gorm:"not null" json:"title"`
	Slug            string    `gorm:"uniqueIndex;not null" json:"slug"`
	Division        string    `gorm:"not null" json:"division"`
	City            string    `gorm:"not null" json:"city"`
	ExperienceLevel string    `gorm:"not null" json:"experience_level"`
	EducationLevel  string    `gorm:"not null" json:"education_level"`
	EmploymentType  string    `gorm:"not null;default:'Fulltime'" json:"employment_type"`
	Description     string    `gorm:"type:text;not null" json:"description"`
	Requirements    string    `gorm:"type:text;not null" json:"requirements"`
	IsActive        *bool     `gorm:"default:true;index" json:"is_active"`
	PublishedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"published_at"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
