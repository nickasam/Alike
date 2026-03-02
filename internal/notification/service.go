package notification

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

func (s *Service) Create(userID, notificationType, title, content string, data map[string]interface{}) error {
	notification := &domain.Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Content: content,
		Data:    data,
		IsRead:  false,
	}
	
	if err := s.db.Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	
	return nil
}

func (s *Service) GetNotifications(userID string, limit, offset int) ([]domain.Notification, error) {
	var notifications []domain.Notification
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
		
	if err != nil {
		return nil, err
	}
	
	return notifications, nil
}

func (s *Service) MarkAsRead(notificationID string) error {
	return s.db.Model(&domain.Notification{}).
		Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

func (s *Service) MarkAllAsRead(userID string) error {
	return s.db.Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}
