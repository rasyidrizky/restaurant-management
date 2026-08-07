package careers

import (
	"errors"

	"project-2026-06-misoastory-be-go/internal/common/dto"
	"project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/utils"
	careerdto "project-2026-06-misoastory-be-go/internal/modules/careers/dto"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

var (
	ErrJobPostNotFound = errors.New("job post not found")
	ErrJobPostConflict = errors.New("job post with this title already exists")
)

// CareerService handles the core business logic and database interactions for Careers.
type CareerService struct {
	db *gorm.DB
}

// NewCareerService creates a new CareerService.
func NewCareerService(db *gorm.DB) *CareerService {
	return &CareerService{db: db}
}

func (s *CareerService) CreateJobPost(req *careerdto.CreateJobPostRequest) (*models.JobPost, error) {
	slug := utils.ToSlug(req.Title)

	var existing models.JobPost
	if err := s.db.Where("slug = ?", slug).First(&existing).Error; err == nil {
		return nil, ErrJobPostConflict
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var jobPost models.JobPost
	copier.Copy(&jobPost, req)
	jobPost.Slug = slug

	if err := s.db.Create(&jobPost).Error; err != nil {
		return nil, err
	}

	return &jobPost, nil
}

func (s *CareerService) GetJobPosts(q *careerdto.JobPostQuery) ([]models.JobPost, dto.Meta, error) {
	var jobPosts []models.JobPost
	var total int64

	query := s.db.Model(&models.JobPost{})

	if q.Search != "" {
		searchTerm := "%" + q.Search + "%"
		query = query.Where("title ILIKE ? OR division ILIKE ? OR city ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, dto.Meta{}, err
	}

	sort := q.Sort
	if sort == "" {
		sort = "created_at desc"
	}

	err := query.Order(sort).Scopes(utils.Paginate(q.Page, q.Limit)).Find(&jobPosts).Error
	meta := utils.CalculateMeta(total, q.Page, q.Limit)

	return jobPosts, meta, err
}

func (s *CareerService) GetJobPostBySlug(slug string) (*models.JobPost, error) {
	var jobPost models.JobPost
	if err := s.db.Where("slug = ?", slug).First(&jobPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJobPostNotFound
		}
		return nil, err
	}
	return &jobPost, nil
}

func (s *CareerService) GetJobPostByID(id uint) (*models.JobPost, error) {
	var jobPost models.JobPost
	if err := s.db.First(&jobPost, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrJobPostNotFound
		}
		return nil, err
	}
	return &jobPost, nil
}

func (s *CareerService) UpdateJobPost(id uint, req *careerdto.UpdateJobPostRequest) (*models.JobPost, error) {
	jobPost, err := s.GetJobPostByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil && *req.Title != jobPost.Title {
		newSlug := utils.ToSlug(*req.Title)
		var existing models.JobPost
		if err := s.db.Where("slug = ? AND id != ?", newSlug, id).First(&existing).Error; err == nil {
			return nil, ErrJobPostConflict
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		jobPost.Slug = newSlug
	}

	if err := copier.CopyWithOption(jobPost, req, copier.Option{IgnoreEmpty: true}); err != nil {
		return nil, err
	}

	if err := s.db.Save(jobPost).Error; err != nil {
		return nil, err
	}

	return jobPost, nil
}

func (s *CareerService) DeleteJobPost(id uint) error {
	jobPost, err := s.GetJobPostByID(id)
	if err != nil {
		return err
	}
	return s.db.Delete(jobPost).Error
}
