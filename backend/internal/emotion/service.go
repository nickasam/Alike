package emotion

import "context"

// Service 是对 Repository 的薄封装，供 WebSocket Hub 广播情绪看板更新。
// 通过接口解耦，使 ws 传输层无需直接依赖 emotion 业务细节。
type Service struct {
	repo *Repository
}

// NewService 创建情绪服务适配器。
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// EmotionBoard 返回频道"今日"情绪看板（供实时 emotion_update 推送）。
// 返回 any 以匹配 ws.EmotionBoardProvider 接口，避免 ws 依赖本包类型。
func (s *Service) EmotionBoard(ctx context.Context, channelID int64) (any, error) {
	return s.repo.BoardByChannel(ctx, channelID, true)
}
