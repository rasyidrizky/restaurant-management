package models

import (
	"time"

	"gorm.io/gorm"
)

type Position struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Description *string        `json:"description,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	Users               []User               `gorm:"foreignKey:PositionID" json:"users,omitempty"`
	PositionPermissions []PositionPermission `gorm:"foreignKey:PositionID" json:"position_permissions,omitempty"`
}
