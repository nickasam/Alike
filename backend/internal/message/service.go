package message

import "context"

// Service 是对 Repository 的薄封装，供 WebSocket Hub 通过 send_message 落库消息。
// 方法签名与 ws.MsgService 接口匹配（结构化实现，无需相互导入）。
type Service struct {
	repo *Repository
}

// NewService 创建消息服务适配器。
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// IsMember 报告用户是否为频道成员。
func (s *Service) IsMember(ctx context.Context, channelID, userID int64) (bool, error) {
	return s.repo.IsMember(ctx, channelID, userID)
}

// CreateMessage 落库一条主消息并返回脱敏后的消息对象（供 WebSocket 广播）。
func (s *Service) CreateMessage(ctx context.Context, channelID, userID int64, content, emotion string, anonymous bool) (any, error) {
	msg, err := s.repo.Create(ctx, channelID, nil, userID, CreateRequest{
		Content:     content,
		Emotion:     emotion,
		IsAnonymous: anonymous,
	})
	if err != nil {
		return nil, err
	}
	msg.mask()
	return msg, nil
}
