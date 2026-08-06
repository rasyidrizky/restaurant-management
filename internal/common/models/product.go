// Package models contains the GORM database models for the application.
package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Slug         string         `gorm:"unique;not null" json:"slug"`
	Description  *string        `json:"description"`
	Price        float64        `gorm:"type:decimal(12,2);not null" json:"price"`
	ImageURL     *string        `gorm:"column:image_url" json:"imageUrl"`
	IsAvailable  bool           `gorm:"column:is_available;default:true" json:"isAvailable"`
	IsBestSeller bool           `gorm:"column:is_best_seller;default:false" json:"isBestSeller"`
	DisplayOrder int            `gorm:"column:display_order;default:0" json:"displayOrder"`
	CategoryID   uint           `gorm:"column:category_id;not null" json:"categoryId"`
	CreatedAt    time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Relations
	Category  *Category         `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category,omitempty"`
	Locations []ProductLocation `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"locations,omitempty"`
}
