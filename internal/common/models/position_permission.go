package models

import (
	"time"
)

type PositionPermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PositionID   uint      `gorm:"not null;uniqueIndex:idx_position_permission" json:"position_id"`
	PermissionID uint      `gorm:"not null;uniqueIndex:idx_position_permission" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`

	Position   Position   `gorm:"foreignKey:PositionID;constraint:OnDelete:CASCADE" json:"position,omitempty"`
	Permission Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE" json:"permission,omitempty"`
}
