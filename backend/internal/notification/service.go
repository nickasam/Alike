package notification

import "context"

// Pusher 抽象向指定用户实时推送通知的能力，由 ws.Hub 实现（可为 nil）。
type Pusher interface {
	NotifyUser(userID int64, payload any)
}

// Service 组合「写库 + 实时推送」：写入通知后经 Pusher 向该用户在线连接推送。
// 实现与 Repository.Create 相同的签名，可直接注入 empathy/message 的 Notifier 接口。
type Service struct {
	repo   *Repository
	pusher Pusher
}

// NewService 创建通知服务。pusher 可为 nil（仅写库，不实时推送）。
func NewService(repo *Repository, pusher Pusher) *Service {
	return &Service{repo: repo, pusher: pusher}
}

// Create 写入一条通知并尝试实时推送（fire-and-forget）。
// 写库失败直接返回错误；推送为尽力而为，用户离线时无接收者。
func (s *Service) Create(ctx context.Context, userID int64, typ, content string, refID *int64) error {
	if err := s.repo.Create(ctx, userID, typ, content, refID); err != nil {
		return err
	}
	if s.pusher != nil {
		s.pusher.NotifyUser(userID, map[string]any{
			"type":    typ,
			"content": content,
			"ref_id":  refID,
		})
	}
	return nil
}
