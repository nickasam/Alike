package message

import "time"

// deletedPlaceholder 是软删除消息对外返回的占位内容。
const deletedPlaceholder = "该消息已删除"

// Author 是消息作者的公开信息。匿名消息不返回该字段。
type Author struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

// Message 是 messages 表的领域模型（含作者信息与匿名/软删除标记）。
type Message struct {
	ID           int64      `json:"id"`
	ChannelID    int64      `json:"channel_id"`
	ParentID     *int64     `json:"parent_id"`
	Content      string     `json:"content"`
	Emotion      string     `json:"emotion,omitempty"`
	IsAnonymous  bool       `json:"is_anonymous"`
	EmpathyCount int        `json:"empathy_count"`
	ReplyCount   int        `json:"reply_count"`
	// Empathized 表示请求方（当前登录用户）是否已对本消息共情。
	Empathized bool       `json:"empathized"`
	IsDeleted    bool       `json:"is_deleted"`
	Author       *Author    `json:"author,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`

	// ClientMsgID 是发送者的客户端幂等标识，仅经 WebSocket 发送时回显，
	// 供发送端将服务端回环消息与本地乐观条目精确合并去重。不入库、不 scan。
	ClientMsgID string `json:"client_msg_id,omitempty"`

	// authorID 供仓储内部做归属判断，不序列化对外。
	authorID int64
}

// mask 依据匿名与软删除标记，隐藏对外不应暴露的字段。
// 匿名消息不返回作者；软删除消息内容替换为占位符并清空作者与情绪。
func (m *Message) mask() {
	if m.IsDeleted {
		m.Content = deletedPlaceholder
		m.Emotion = ""
		m.Author = nil
		return
	}
	if m.IsAnonymous {
		m.Author = nil
	}
}

// CreateRequest 是发布消息 / 线程回复的请求体。
type CreateRequest struct {
	Content     string `json:"content" binding:"required,min=1,max=5000"`
	Emotion     string `json:"emotion" binding:"max=50"`
	IsAnonymous bool   `json:"is_anonymous"`
	// ClientMsgID 是客户端生成的幂等去重标识（UUID），可选。
	ClientMsgID string `json:"client_msg_id" binding:"max=64"`
}
