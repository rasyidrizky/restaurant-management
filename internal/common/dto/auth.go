package dto

import (
	"project-2026-06-misoastory-be-go/internal/common/models"

	"github.com/jinzhu/copier"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token    string       `json:"token"`
	User     UserResponse `json:"user"`
}

type UserResponse struct {
	ID         uint   `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	PositionID uint   `json:"position_id"`
}

func MapToUserResponse(user *models.User) UserResponse {
	var resp UserResponse
	copier.Copy(&resp, user)
	return resp
}

func MapToUserResponses(users []models.User) []UserResponse {
	var responses []UserResponse
	copier.Copy(&responses, &users)
	if responses == nil {
		return []UserResponse{}
	}
	return responses
}
