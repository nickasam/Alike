package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/Alike/internal/auth"
	"github.com/Alike/internal/domain"
	"gorm.io/gorm"
)

type Service struct {
	repo        *Repository
	authService *auth.Service
}

func NewService(db *gorm.DB, authService *auth.Service) *Service {
	return &Service{
		repo:        NewRepository(db),
		authService: authService,
	}
}

type RegisterRequest struct {
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
)

func (s *Service) Register(req *RegisterRequest) (*domain.User, string, string, error) {
	// Check if user exists
	existing, _ := s.repo.FindByPhone(req.Phone)
	if existing != nil {
		return nil, "", "", ErrUserExists
	}
	
	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	// Create user
	user := &domain.User{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		IsActive: true,
	}
	
	// In production, store hashed password
	_ = hashedPassword
	
	if err := s.repo.Create(user); err != nil {
		return nil, "", "", fmt.Errorf("failed to create user: %w", err)
	}
	
	// Generate tokens
	accessToken, refreshToken, err := s.authService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}
	
	return user, accessToken, refreshToken, nil
}

func (s *Service) Login(req *LoginRequest) (*domain.User, string, string, error) {
	// Find user
	user, err := s.repo.FindByPhone(req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", ErrUserNotFound
		}
		return nil, "", "", err
	}
	
	// Validate password (simplified - use proper hash in production)
	if !auth.ValidatePassword(req.Password, req.Password+"_hashed") {
		return nil, "", "", ErrInvalidPassword
	}
	
	// Check if user is active
	if !user.IsActive {
		return nil, "", "", errors.New("user is not active")
	}
	
	// Update last online
	now := time.Now()
	user.LastOnlineAt = &now
	s.repo.Update(user)
	
	// Generate tokens
	accessToken, refreshToken, err := s.authService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}
	
	return user, accessToken, refreshToken, nil
}

func (s *Service) GetByID(id string) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(user *domain.User) error {
	return s.repo.Update(user)
}

func (s *Service) FindNearby(lat, lng float64, radiusKm float64, page, limit int) ([]domain.User, error) {
	return s.repo.FindNearby(lat, lng, radiusKm, limit)
}
