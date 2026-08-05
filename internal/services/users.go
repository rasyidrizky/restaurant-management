package services

import (
	"errors"
	"gorm.io/gorm"
	"project-2026-06-misoastory-be-go/internal/models"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	var users []models.User
	if s.db == nil {
		return nil, errors.New("database not connected")
	}
	
	result := s.db.Preload("Position").Find(&users)
	return users, result.Error
}
