package services

import (
	"errors"
	"strings"

	"project-2026-06-misoastory-be-go/internal/dto"
	"project-2026-06-misoastory-be-go/internal/models"
	"project-2026-06-misoastory-be-go/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryConflict = errors.New("category slug already exists")
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		db: db,
	}
}

func (s *CategoryService) GetAllCategories(search string) ([]models.Category, error) {
	var categories []models.Category
	query := s.db.Model(&models.Category{})

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	err := query.Order("display_order asc, name asc").Find(&categories).Error
	return categories, err
}

func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) CreateCategory(req *dto.CreateCategoryRequest) (*models.Category, error) {
	slug := utils.ToSlug(req.Name)

	// Check if slug exists
	var existingCategory models.Category
	if err := s.db.Where("slug = ?", slug).First(&existingCategory).Error; err == nil {
		return nil, ErrCategoryConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := models.Category{
		Name: req.Name,
		Slug: slug,
	}

	if req.Description != nil {
		category.Description = req.Description
	}
	if req.DisplayOrder != nil {
		category.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	} else {
		category.IsActive = true // default
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		category.Name = *req.Name
		newSlug := utils.ToSlug(*req.Name)
		
		// Only check for conflict if the slug is actually changing
		if newSlug != category.Slug {
			var existingCategory models.Category
			if err := s.db.Where("slug = ?", newSlug).First(&existingCategory).Error; err == nil {
				return nil, ErrCategoryConflict
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			category.Slug = newSlug
		}
	}

	if req.Description != nil {
		category.Description = req.Description
	}
	if req.DisplayOrder != nil {
		category.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.db.Save(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) DeleteCategory(id uint) error {
	category, err := s.GetCategoryByID(id)
	if err != nil {
		return err
	}

	return s.db.Delete(category).Error
}
