package locations

import (
	"errors"
	"project-2026-06-misoastory-be-go/internal/common/dto"
	locationdto "project-2026-06-misoastory-be-go/internal/modules/locations/dto"
	"project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/utils"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var (
	ErrLocationNotFound = errors.New("location not found")
	ErrLocationConflict = errors.New("location with slug already exists")
)

// LocationService handles business logic for Locations
type LocationService struct {
	db *gorm.DB
}

// NewLocationService acts as the constructor for LocationService
func NewLocationService(db *gorm.DB) *LocationService {
	return &LocationService{
		db: db,
	}
}

func (s *LocationService) GetLocations(req *locationdto.LocationQuery) ([]models.Location, dto.Meta, error) {
	var locations []models.Location
	var total int64

	query := s.db.Model(&models.Location{})
	
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR city ILIKE ?", searchTerm, searchTerm)
	}

	// Count total before paginating
	if err := query.Count(&total).Error; err != nil {
		return nil, dto.Meta{}, err
	}

	// Default sort if not provided
	sort := req.Sort
	if sort == "" {
		sort = "name asc"
	}

	err := query.Order(sort).Scopes(utils.Paginate(req.Page, req.Limit)).Find(&locations).Error
	meta := utils.CalculateMeta(total, req.Page, req.Limit)
	
	return locations, meta, err
}

func (s *LocationService) GetLocationByID(id uint) (*models.Location, error) {
	var location models.Location
	if err := s.db.First(&location, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return &location, nil
}

func (s *LocationService) CreateLocation(req *locationdto.CreateLocationRequest) (*models.Location, error) {
	slug := utils.ToSlug(req.Name)

	// Check conflict
	var existing models.Location
	if err := s.db.Where("slug = ?", slug).First(&existing).Error; err == nil {
		return nil, ErrLocationConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	location := models.Location{
		Slug: slug,
	}
	copier.Copy(&location, req)

	if err := s.db.Create(&location).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (s *LocationService) UpdateLocation(id uint, req *locationdto.UpdateLocationRequest) (*models.Location, error) {
	location, err := s.GetLocationByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil && *req.Name != location.Name {
		newSlug := utils.ToSlug(*req.Name)
		var existing models.Location
		if err := s.db.Where("slug = ? AND id != ?", newSlug, id).First(&existing).Error; err == nil {
			return nil, ErrLocationConflict
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		location.Slug = newSlug
	}

	// Automatically copy all non-nil fields
	copier.CopyWithOption(location, req, copier.Option{IgnoreEmpty: true})

	if err := s.db.Save(location).Error; err != nil {
		return nil, err
	}

	return location, nil
}

func (s *LocationService) DeleteLocation(id uint) error {
	location, err := s.GetLocationByID(id)
	if err != nil {
		return err
	}
	
	return s.db.Delete(location).Error
}
