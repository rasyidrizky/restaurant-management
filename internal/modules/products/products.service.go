package products

import (
	"errors"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	"project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/utils"
	producttypes "project-2026-06-misoastory-be-go/internal/modules/products/types"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductConflict = errors.New("product with this name already exists")
)

// ProductService handles the core business logic and database interactions for Products.
type ProductService struct {
	db *gorm.DB
}

// NewProductService acts as the constructor for ProductService, injecting the GORM DB instance.
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) CreateProduct(req *producttypes.CreateProductRequest) (*models.Product, error) {
	slug := utils.ToSlug(req.Name)

	var existing models.Product
	if err := s.db.Where("slug = ?", slug).First(&existing).Error; err == nil {
		return nil, ErrProductConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var product models.Product
	if err := copier.Copy(&product, req); err != nil {
		return nil, err
	}
	product.Slug = slug

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		if len(req.LocationIDs) > 0 {
			var productLocations []models.ProductLocation
			for _, locID := range req.LocationIDs {
				productLocations = append(productLocations, models.ProductLocation{
					ProductID:   product.ID,
					LocationID:  locID,
					IsAvailable: true,
				})
			}
			if err := tx.Create(&productLocations).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetProductByID(product.ID)
}

func (s *ProductService) UpdateProduct(id uint, req *producttypes.UpdateProductRequest) (*models.Product, error) {
	product, err := s.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil && *req.Name != product.Name {
		newSlug := utils.ToSlug(*req.Name)
		var existing models.Product
		if err := s.db.Where("slug = ? AND id != ?", newSlug, id).First(&existing).Error; err == nil {
			return nil, ErrProductConflict
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		product.Slug = newSlug
	}

	if err := copier.CopyWithOption(product, req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		if req.LocationIDs != nil {
			// Replace locations mapping
			if err := tx.Where("product_id = ?", id).Delete(&models.ProductLocation{}).Error; err != nil {
				return err
			}
			if len(req.LocationIDs) > 0 {
				var productLocations []models.ProductLocation
				for _, locID := range req.LocationIDs {
					productLocations = append(productLocations, models.ProductLocation{
						ProductID:   id,
						LocationID:  locID,
						IsAvailable: true,
					})
				}
				if err := tx.Create(&productLocations).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetProductByID(product.ID)
}

func (s *ProductService) GetProducts(q *producttypes.ProductQuery) ([]models.Product, dto.Meta, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{}).
		Preload("Category").
		Preload("Locations")

	if q.Search != "" {
		query = query.Where("name ILIKE ?", "%"+q.Search+"%")
	}
	if q.CategoryID != nil {
		query = query.Where("category_id = ?", *q.CategoryID)
	}
	if q.LocationID != nil {
		query = query.Joins("JOIN product_locations ON product_locations.product_id = products.id").
			Where("product_locations.location_id = ?", *q.LocationID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, dto.Meta{}, err
	}

	if q.Sort != "" {
		query = query.Order(q.Sort)
	} else {
		query = query.Order("display_order asc, id desc")
	}

	if err := query.Scopes(utils.Paginate(q.Page, q.Limit)).Find(&products).Error; err != nil {
		return nil, dto.Meta{}, err
	}

	meta := utils.CalculateMeta(total, q.Page, q.Limit)
	return products, meta, nil
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := s.db.Preload("Category").Preload("Locations").Preload("Locations.Location").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) DeleteProduct(id uint) error {
	product, err := s.GetProductByID(id)
	if err != nil {
		return err
	}
	return s.db.Delete(product).Error
}
