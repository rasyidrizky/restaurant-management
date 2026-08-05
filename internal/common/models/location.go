package models

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Name                string         `gorm:"not null" json:"name"`
	Slug                string         `gorm:"unique;not null" json:"slug"`
	Address             string         `gorm:"not null" json:"address"`
	City                string         `gorm:"not null" json:"city"`
	Phone               *string        `json:"phone,omitempty"`
	Latitude            *float64       `json:"latitude,omitempty"`
	Longitude           *float64       `json:"longitude,omitempty"`
	OpeningTime         *string        `json:"opening_time,omitempty"`
	ClosingTime         *string        `json:"closing_time,omitempty"`
	IsActive            bool           `gorm:"default:true" json:"is_active"`
	IsOpen24Hours       bool           `gorm:"default:false" json:"is_open_24_hours"`
	HasDineIn           bool           `gorm:"default:true" json:"has_dine_in"`
	SupportsHomeService bool           `gorm:"default:false" json:"supports_home_service"`
	SupportsEvents      bool           `gorm:"default:false" json:"supports_events"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}
