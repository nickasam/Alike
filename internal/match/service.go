package match

import (
	"fmt"

	"github.com/Alike/internal/domain"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateLike(likerID, likedID string) error {
	like := &domain.Like{
		LikerID: likerID,
		LikedID: likedID,
	}
	
	if err := s.db.Create(like).Error; err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}
	
	// Check for mutual like
	var existingLike domain.Like
	err := s.db.Where("liker_id = ? AND liked_id = ?", likedID, likerID).First(&existingLike).Error
	if err == nil {
		// Mutual like detected - create match
		return s.CreateMatch(likerID, likedID)
	}
	
	return nil
}

func (s *Service) CreateMatch(user1ID, user2ID string) error {
	match := &domain.Match{
		User1ID:  user1ID,
		User2ID:  user2ID,
		IsActive: true,
	}
	
	if err := s.db.Create(match).Error; err != nil {
		return fmt.Errorf("failed to create match: %w", err)
	}
	
	return nil
}

func (s *Service) GetMatches(userID string) ([]domain.Match, error) {
	var matches []domain.Match
	err := s.db.Where("user1_id = ? OR user2_id = ?", userID, userID).
		Where("is_active = ?", true).
		Order("last_message_at DESC").
		Find(&matches).Error
		
	if err != nil {
		return nil, err
	}
	
	return matches, nil
}

func (s *Service) GetMatch(id string) (*domain.Match, error) {
	var match domain.Match
	err := s.db.Where("id = ?", id).First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}
