package services

import (
	"errors"
	"project-2026-06-misoastory-be-go/internal/database"
	"project-2026-06-misoastory-be-go/internal/dto"
	"project-2026-06-misoastory-be-go/internal/models"
	"project-2026-06-misoastory-be-go/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrLocationNotFound = errors.New("location not found")
	ErrLocationConflict = errors.New("location with slug already exists")
)

type LocationService struct{}

func NewLocationService() *LocationService {
	return &LocationService{}
}

func (s *LocationService) GetAllLocations(search string) ([]models.Location, error) {
	var locations []models.Location
	if database.DB == nil {
		return nil, errors.New("database not connected")
	}

	query := database.DB.Model(&models.Location{})
	
	if search != "" {
		query = query.Where("name ILIKE ? OR city ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&locations)
	return locations, result.Error
}

func (s *LocationService) GetLocationByID(id uint) (*models.Location, error) {
	var location models.Location
	if err := database.DB.First(&location, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return &location, nil
}

func (s *LocationService) CreateLocation(req *dto.CreateLocationRequest) (*models.Location, error) {
	slug := utils.ToSlug(req.Name)

	// Check conflict
	var count int64
	database.DB.Model(&models.Location{}).Where("slug = ?", slug).Count(&count)
	if count > 0 {
		return nil, ErrLocationConflict
	}

	location := models.Location{
		Name:                req.Name,
		Slug:                slug,
		Address:             req.Address,
		City:                req.City,
		Phone:               req.Phone,
		Latitude:            req.Latitude,
		Longitude:           req.Longitude,
		OpeningTime:         req.OpeningTime,
		ClosingTime:         req.ClosingTime,
		IsActive:            req.IsActive,
		IsOpen24Hours:       req.IsOpen24Hours,
		HasDineIn:           req.HasDineIn,
		SupportsHomeService: req.SupportsHomeService,
		SupportsEvents:      req.SupportsEvents,
	}

	if err := database.DB.Create(&location).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (s *LocationService) UpdateLocation(id uint, req *dto.UpdateLocationRequest) (*models.Location, error) {
	location, err := s.GetLocationByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil && *req.Name != location.Name {
		newSlug := utils.ToSlug(*req.Name)
		var count int64
		database.DB.Model(&models.Location{}).Where("slug = ? AND id != ?", newSlug, id).Count(&count)
		if count > 0 {
			return nil, ErrLocationConflict
		}
		location.Name = *req.Name
		location.Slug = newSlug
	}

	if req.Address != nil {
		location.Address = *req.Address
	}
	if req.City != nil {
		location.City = *req.City
	}
	if req.Phone != nil {
		location.Phone = req.Phone
	}
	if req.Latitude != nil {
		location.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		location.Longitude = req.Longitude
	}
	if req.OpeningTime != nil {
		location.OpeningTime = req.OpeningTime
	}
	if req.ClosingTime != nil {
		location.ClosingTime = req.ClosingTime
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}
	if req.IsOpen24Hours != nil {
		location.IsOpen24Hours = *req.IsOpen24Hours
	}
	if req.HasDineIn != nil {
		location.HasDineIn = *req.HasDineIn
	}
	if req.SupportsHomeService != nil {
		location.SupportsHomeService = *req.SupportsHomeService
	}
	if req.SupportsEvents != nil {
		location.SupportsEvents = *req.SupportsEvents
	}

	if err := database.DB.Save(location).Error; err != nil {
		return nil, err
	}

	return location, nil
}

func (s *LocationService) DeleteLocation(id uint) error {
	location, err := s.GetLocationByID(id)
	if err != nil {
		return err
	}
	
	return database.DB.Delete(location).Error
}
