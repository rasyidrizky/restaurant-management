package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`
	Description *string        `json:"description,omitempty"`
	Resource    string         `gorm:"not null" json:"resource"` // e.g., 'USER', 'TRANSACTION'
	Action      string         `gorm:"not null" json:"action"`   // e.g., 'VIEW', 'ADD'
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	PositionPermissions []PositionPermission `gorm:"foreignKey:PermissionID" json:"-"`
}
