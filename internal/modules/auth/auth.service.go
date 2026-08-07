package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	authdto "project-2026-06-misoastory-be-go/internal/modules/auth/dto"
	userdto "project-2026-06-misoastory-be-go/internal/modules/users/dto"
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

func (s *AuthService) Register(req *authdto.RegisterRequest) (*authdto.AuthResponse, error) {
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

			// Seed basic admin permissions automatically for implemented modules
			permissions := []models.Permission{
				// User Management
				{Name: "VIEW_USER", Resource: "USER", Action: "VIEW"},
				{Name: "ADD_USER", Resource: "USER", Action: "ADD"},
				{Name: "UPDATE_USER", Resource: "USER", Action: "UPDATE"},
				{Name: "DELETE_USER", Resource: "USER", Action: "DELETE"},
				{Name: "MANAGE_USER_PERMISSION", Resource: "USER", Action: "MANAGE_PERMISSION"},
				{Name: "CHANGE_USER_POSITION", Resource: "USER", Action: "CHANGE_POSITION"},

				// Position Management
				{Name: "VIEW_POSITION", Resource: "POSITION", Action: "VIEW"},
				{Name: "ADD_POSITION", Resource: "POSITION", Action: "ADD"},
				{Name: "UPDATE_POSITION", Resource: "POSITION", Action: "UPDATE"},
				{Name: "DELETE_POSITION", Resource: "POSITION", Action: "DELETE"},

				// Permission Management
				{Name: "VIEW_PERMISSION", Resource: "PERMISSION", Action: "VIEW"},
				{Name: "ADD_PERMISSION", Resource: "PERMISSION", Action: "ADD"},
				{Name: "UPDATE_PERMISSION", Resource: "PERMISSION", Action: "UPDATE"},
				{Name: "DELETE_PERMISSION", Resource: "PERMISSION", Action: "DELETE"},

				// Category Management
				{Name: "VIEW_CATEGORY", Resource: "CATEGORY", Action: "VIEW"},
				{Name: "ADD_CATEGORY", Resource: "CATEGORY", Action: "ADD"},
				{Name: "UPDATE_CATEGORY", Resource: "CATEGORY", Action: "UPDATE"},
				{Name: "DELETE_CATEGORY", Resource: "CATEGORY", Action: "DELETE"},

				// Product Management
				{Name: "VIEW_PRODUCT", Resource: "PRODUCT", Action: "VIEW"},
				{Name: "ADD_PRODUCT", Resource: "PRODUCT", Action: "ADD"},
				{Name: "UPDATE_PRODUCT", Resource: "PRODUCT", Action: "UPDATE"},
				{Name: "DELETE_PRODUCT", Resource: "PRODUCT", Action: "DELETE"},

				// Location Management
				{Name: "VIEW_LOCATION", Resource: "LOCATION", Action: "VIEW"},
				{Name: "ADD_LOCATION", Resource: "LOCATION", Action: "ADD"},
				{Name: "UPDATE_LOCATION", Resource: "LOCATION", Action: "UPDATE"},
				{Name: "DELETE_LOCATION", Resource: "LOCATION", Action: "DELETE"},
			}
			
			for _, p := range permissions {
				var perm models.Permission
				s.db.Where("resource = ? AND action = ?", p.Resource, p.Action).FirstOrCreate(&perm, p)
				
				s.db.Create(&models.PositionPermission{
					PositionID:   position.ID,
					PermissionID: perm.ID,
				})
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

	return &authdto.AuthResponse{
		Token: token,
		User:  userdto.MapToUserResponse(&user),
	}, nil
}

func (s *AuthService) Login(req *authdto.LoginRequest) (*authdto.AuthResponse, error) {
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

	return &authdto.AuthResponse{
		Token: token,
		User:  userdto.MapToUserResponse(&user),
	}, nil
}
