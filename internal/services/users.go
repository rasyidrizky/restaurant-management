package services

import (
	"errors"
	"project-2026-06-misoastory-be-go/internal/database"
	"project-2026-06-misoastory-be-go/internal/models"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	var users []models.User
	if database.DB == nil {
		return nil, errors.New("database not connected")
	}
	
	result := database.DB.Preload("Position").Find(&users)
	return users, result.Error
}
