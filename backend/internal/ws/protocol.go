package ws

import "encoding/json"

// 客户端 → 服务端事件类型。
const (
	EventAuth         = "auth"          // 首帧鉴权：{ token }
	EventJoinChannel  = "join_channel"  // 加入频道：{ channel_id }
	EventLeaveChannel = "leave_channel" // 离开频道：{ channel_id }
	EventTyping       = "typing"        // 正在输入：{ channel_id }
	EventSendMessage  = "send_message"  // 发送消息：{ channel_id, content, emotion, is_anonymous, client_msg_id }
	EventPong         = "pong"          // 心跳响应
)

// 服务端 → 客户端事件类型。
const (
	EventAuthOK       = "auth_ok"       // 鉴权成功
	EventNewMessage   = "new_message"   // 频道新消息
	EventThreadReply  = "thread_reply"  // 线程新回复
	EventMessageDeleted = "message_deleted" // 消息被软删除
	EventEmpathy      = "empathy"       // 共情变更
	EventUserJoined   = "user_joined"   // 有用户加入频道
	EventEmotionUpdate = "emotion_update" // 情绪看板更新
	EventNotification = "notification"  // 通知
	EventError        = "error"         // 错误
	EventPing         = "ping"          // 心跳
)

// Envelope 是所有 WebSocket 帧的统一信封。
// 入站帧的业务字段放在 Data（原始 JSON）中按类型解析；
// 出站帧直接填充 Data 为可序列化对象。
// ChannelID 用于按频道路由；UserID 用于按用户定向路由（如通知），二者互斥使用。
type Envelope struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data,omitempty"`
	ChannelID int64           `json:"channel_id,omitempty"`
	UserID    int64           `json:"user_id,omitempty"`
}

// outbound 构造一个带对象 Data 的出站信封（Data 会被序列化）。
func outbound(eventType string, channelID int64, data any) Envelope {
	raw, _ := json.Marshal(data)
	return Envelope{Type: eventType, ChannelID: channelID, Data: raw}
}

// outboundUser 构造一个按用户定向的出站信封（用于通知等点对点推送）。
func outboundUser(eventType string, userID int64, data any) Envelope {
	raw, _ := json.Marshal(data)
	return Envelope{Type: eventType, UserID: userID, Data: raw}
}

// authData 是首帧鉴权帧的业务体。
type authData struct {
	Token string `json:"token"`
}

// channelData 是携带频道 ID 的通用入站业务体（join/leave/typing）。
type channelData struct {
	ChannelID int64 `json:"channel_id"`
}

// sendMessageData 是 send_message 事件的业务体。
type sendMessageData struct {
	ChannelID   int64  `json:"channel_id"`
	Content     string `json:"content"`
	Emotion     string `json:"emotion"`
	IsAnonymous bool   `json:"is_anonymous"`
	ClientMsgID string `json:"client_msg_id"`
}
