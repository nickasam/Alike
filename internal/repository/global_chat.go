package repository

import (
	"github.com/Alike/internal/domain"
	"gorm.io/gorm"
)

type GlobalChatRepository struct {
	db *gorm.DB
}

func NewGlobalChatRepository(db *gorm.DB) *GlobalChatRepository {
	return &GlobalChatRepository{db: db}
}

// CreateRoom 创建聊天室
func (r *GlobalChatRepository) CreateRoom(room *domain.GlobalChatRoom) error {
	return r.db.Create(room).Error
}

// GetRoom 获取聊天室
func (r *GlobalChatRepository) GetRoom(id string) (*domain.GlobalChatRoom, error) {
	var room domain.GlobalChatRoom
	err := r.db.Where("id = ?", id).First(&room).Error
	return &room, err
}

// CreateMessage 创建消息
func (r *GlobalChatRepository) CreateMessage(msg *domain.GlobalMessage) error {
	return r.db.Create(msg).Error
}

// GetMessages 获取消息列表
func (r *GlobalChatRepository) GetMessages(roomID string, limit int) ([]domain.GlobalMessage, error) {
	var messages []domain.GlobalMessage
	err := r.db.Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// JoinRoom 加入聊天室
func (r *GlobalChatRepository) JoinRoom(userID, roomID string) error {
	// 所有用户都可以访问全局聊天室
	return nil
}
