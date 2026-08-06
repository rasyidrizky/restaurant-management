package users

import (
	"errors"
	"gorm.io/gorm"
	"project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/dto"
	"project-2026-06-misoastory-be-go/internal/common/utils"
	usertypes "project-2026-06-misoastory-be-go/internal/modules/users/types"
)

// UserService handles business logic for Users
type UserService struct {
	db *gorm.DB
}

// NewUserService acts as the constructor for UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetUsers(req *usertypes.GetAllUsersRequest) ([]models.User, dto.Meta, error) {
	var users []models.User
	var total int64
	if s.db == nil {
		return nil, dto.Meta{}, errors.New("database not connected")
	}
	
	query := s.db.Model(&models.User{}).Preload("Position")

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Count total before paginating
	if err := query.Count(&total).Error; err != nil {
		return nil, dto.Meta{}, err
	}

	// Default sort if not provided
	sort := req.Sort
	if sort == "" {
		sort = "first_name asc"
	}

	err := query.Order(sort).Scopes(utils.Paginate(req.Page, req.Limit)).Find(&users).Error
	meta := utils.CalculateMeta(total, req.Page, req.Limit)
	
	return users, meta, err
}
