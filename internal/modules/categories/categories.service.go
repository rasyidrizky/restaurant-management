package categories

import (
	"errors"
	"strings"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	categorytypes "project-2026-06-misoastory-be-go/internal/modules/categories/types"
	"project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/utils"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryConflict = errors.New("category slug already exists")
)

// CategoryService handles business logic for Categories
type CategoryService struct {
	db *gorm.DB
}

// NewCategoryService acts as the constructor for CategoryService
func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		db: db,
	}
}

func (s *CategoryService) GetAllCategories(req *categorytypes.GetAllCategoriesRequest) ([]models.Category, dto.Meta, error) {
	var categories []models.Category
	var total int64
	query := s.db.Model(&models.Category{})

	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	// Count total before paginating
	if err := query.Count(&total).Error; err != nil {
		return nil, dto.Meta{}, err
	}

	// Default sort if not provided
	sort := req.Sort
	if sort == "" {
		sort = "display_order asc, name asc"
	}

	err := query.Order(sort).Scopes(utils.Paginate(req.Page, req.Limit)).Find(&categories).Error
	meta := utils.CalculateMeta(total, req.Page, req.Limit)
	
	return categories, meta, err
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

func (s *CategoryService) CreateCategory(req *categorytypes.CreateCategoryRequest) (*models.Category, error) {
	slug := utils.ToSlug(req.Name)

	// Check if slug exists
	var existingCategory models.Category
	if err := s.db.Where("slug = ?", slug).First(&existingCategory).Error; err == nil {
		return nil, ErrCategoryConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := models.Category{
		Slug: slug,
	}
	copier.Copy(&category, req)
	
	if req.IsActive == nil {
		category.IsActive = true // default
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoryService) UpdateCategory(id uint, req *categorytypes.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		newSlug := utils.ToSlug(*req.Name)
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

	// Automatically copy all non-nil fields
	copier.CopyWithOption(category, req, copier.Option{IgnoreEmpty: true})

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
