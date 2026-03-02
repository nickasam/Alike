package chat

import (
	"github.com/Alike/internal/domain"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetChats(userID string) ([]domain.Chat, error) {
	var chats []domain.Chat
	err := s.db.Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("last_message_at DESC").
		Find(&chats).Error
		
	if err != nil {
		return nil, err
	}
	
	return chats, nil
}

func (s *Service) GetChat(id string) (*domain.Chat, error) {
	var chat domain.Chat
	err := s.db.Where("id = ?", id).First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (s *Service) GetMessages(chatID string, page, limit int) ([]domain.Message, error) {
	var messages []domain.Message
	offset := (page - 1) * limit
	
	err := s.db.Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
		
	if err != nil {
		return nil, err
	}
	
	return messages, nil
}

func (s *Service) CreateMessage(message *domain.Message) error {
	return s.db.Create(message).Error
}

func (s *Service) MarkMessageAsRead(messageID string) error {
	return s.db.Model(&domain.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}
