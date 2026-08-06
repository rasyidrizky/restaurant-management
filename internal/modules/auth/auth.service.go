package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	authtypes "project-2026-06-misoastory-be-go/internal/modules/auth/types"
	usertypes "project-2026-06-misoastory-be-go/internal/modules/users/types"
	"project-2026-06-misoastory-be-go/internal/common/models"
	"project-2026-06-misoastory-be-go/internal/common/utils"
)

// AuthService handles business logic for Authentication
type AuthService struct {
	db *gorm.DB
}

// NewAuthService acts as the constructor for AuthService
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(req *authtypes.RegisterRequest) (*authtypes.AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already in use")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Fetch or Create default 'Admin' position for seeding purposes
	var position models.Position
	if err := s.db.Where("name = ?", "Admin").First(&position).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			position = models.Position{Name: "Admin", Description: utils.StringPtr("Administrator")}
			if err := s.db.Create(&position).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   string(hashedPassword),
		PositionID: position.ID,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID, user.PositionID)
	if err != nil {
		return nil, err
	}

	return &authtypes.AuthResponse{
		Token: token,
		User:  usertypes.MapToUserResponse(&user),
	}, nil
}

func (s *AuthService) Login(req *authtypes.LoginRequest) (*authtypes.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.PositionID)
	if err != nil {
		return nil, err
	}

	return &authtypes.AuthResponse{
		Token: token,
		User:  usertypes.MapToUserResponse(&user),
	}, nil
}
