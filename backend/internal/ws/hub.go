package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
)

// MsgService 抽象 Hub 处理 send_message 所需的消息业务能力。
// 由 message.Repository 的适配器实现（在 router 层注入），避免 ws ↔ message 循环依赖。
type MsgService interface {
	// IsMember 报告用户是否为频道成员。
	IsMember(ctx context.Context, channelID, userID int64) (bool, error)
	// CreateMessage 落库一条主消息，返回可序列化的消息对象（已脱敏）。
	CreateMessage(ctx context.Context, channelID, userID int64, content, emotion string, anonymous bool) (any, error)
}

// Hub 管理所有 WebSocket 连接，负责本地广播与 Redis Pub/Sub 跨实例广播。
type Hub struct {
	mu       sync.RWMutex
	clients  map[*Client]struct{}
	channels map[int64]map[*Client]struct{} // channelID -> 订阅该频道的客户端集合

	svc    MsgService
	pubsub *PubSub
}

// NewHub 创建 Hub。svc / pubsub 均可为 nil（对应能力降级）。
func NewHub(svc MsgService, pubsub *PubSub) *Hub {
	h := &Hub{
		clients:  make(map[*Client]struct{}),
		channels: make(map[int64]map[*Client]struct{}),
		svc:      svc,
		pubsub:   pubsub,
	}
	if pubsub != nil {
		pubsub.OnMessage(h.deliverLocal)
		pubsub.Start()
	}
	return h
}

// register 登记一个新客户端。
func (h *Hub) register(c *Client) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
}

// unregister 注销客户端并清理其所有频道订阅与出站队列。
func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	if _, ok := h.clients[c]; !ok {
		h.mu.Unlock()
		return
	}
	delete(h.clients, c)
	for id := range c.channels {
		if set := h.channels[id]; set != nil {
			delete(set, c)
			if len(set) == 0 {
				delete(h.channels, id)
			}
		}
	}
	h.mu.Unlock()
	c.closeSend()
}

// joinChannel 将客户端加入频道的本地订阅集合。
func (h *Hub) joinChannel(c *Client, channelID int64) {
	h.mu.Lock()
	set := h.channels[channelID]
	if set == nil {
		set = make(map[*Client]struct{})
		h.channels[channelID] = set
	}
	set[c] = struct{}{}
	h.mu.Unlock()
	c.subscribe(channelID)
}

// leaveChannel 将客户端移出频道的本地订阅集合。
func (h *Hub) leaveChannel(c *Client, channelID int64) {
	h.mu.Lock()
	if set := h.channels[channelID]; set != nil {
		delete(set, c)
		if len(set) == 0 {
			delete(h.channels, channelID)
		}
	}
	h.mu.Unlock()
	c.unsubscribe(channelID)
}

// deliverLocal 将一个信封投递给本实例中订阅了该频道的所有客户端。
func (h *Hub) deliverLocal(evt Envelope) {
	b, err := json.Marshal(evt)
	if err != nil {
		return
	}
	h.mu.RLock()
	set := h.channels[evt.ChannelID]
	targets := make([]*Client, 0, len(set))
	for c := range set {
		targets = append(targets, c)
	}
	h.mu.RUnlock()

	for _, c := range targets {
		if !c.enqueue(b) {
			// 出站缓冲已满：视为慢客户端，断开以释放资源。
			go h.unregister(c)
		}
	}
}

// publish 将事件广播到频道：优先经 Redis Pub/Sub 跨实例，否则退化为本地广播。
func (h *Hub) publish(evt Envelope) {
	if h.pubsub != nil {
		if err := h.pubsub.Publish(evt); err == nil {
			return
		}
		// Redis 发布失败时降级为本地广播，保证单实例可用。
		log.Printf("[WARN] ws: redis publish failed, fallback to local broadcast")
	}
	h.deliverLocal(evt)
}

// BroadcastNewMessage 实现 message.Broadcaster：广播频道新消息。
func (h *Hub) BroadcastNewMessage(channelID int64, payload any) {
	h.publish(outbound(EventNewMessage, channelID, payload))
}

// BroadcastThreadReply 实现 message.Broadcaster：广播线程新回复。
func (h *Hub) BroadcastThreadReply(channelID, parentID int64, payload any) {
	h.publish(outbound(EventThreadReply, channelID, map[string]any{"parent_id": parentID, "reply": payload}))
}

// handleClientEvent 解析并处理一条客户端入站帧。
func (h *Hub) handleClientEvent(c *Client, raw []byte) {
	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		c.sendEvent(errorEvent("非法的消息格式"))
		return
	}

	switch env.Type {
	case EventPong:
		// 心跳响应，无需处理（readPump 已刷新读超时）。
	case EventJoinChannel:
		h.onJoin(c, env)
	case EventLeaveChannel:
		h.onLeave(c, env)
	case EventTyping:
		h.onTyping(c, env)
	case EventSendMessage:
		h.onSendMessage(c, env)
	default:
		c.sendEvent(errorEvent("未知的事件类型: " + env.Type))
	}
}

func (h *Hub) onJoin(c *Client, env Envelope) {
	var d channelData
	if !decodeData(env, &d) || d.ChannelID <= 0 {
		c.sendEvent(errorEvent("channel_id 无效"))
		return
	}
	if h.svc != nil {
		ok, err := h.svc.IsMember(context.Background(), d.ChannelID, c.userID)
		if err != nil {
			c.sendEvent(errorEvent("加入频道失败"))
			return
		}
		if !ok {
			c.sendEvent(errorEvent("请先加入该频道"))
			return
		}
	}
	h.joinChannel(c, d.ChannelID)
	c.sendEvent(outbound(EventUserJoined, d.ChannelID, map[string]any{"user_id": c.userID}))
}

func (h *Hub) onLeave(c *Client, env Envelope) {
	var d channelData
	if !decodeData(env, &d) || d.ChannelID <= 0 {
		c.sendEvent(errorEvent("channel_id 无效"))
		return
	}
	h.leaveChannel(c, d.ChannelID)
}

func (h *Hub) onTyping(c *Client, env Envelope) {
	var d channelData
	if !decodeData(env, &d) || d.ChannelID <= 0 {
		return
	}
	if !c.isSubscribed(d.ChannelID) {
		return
	}
	// typing 仅在频道内广播，不落库。
	h.publish(outbound(EventTyping, d.ChannelID, map[string]any{"user_id": c.userID}))
}

func (h *Hub) onSendMessage(c *Client, env Envelope) {
	var d sendMessageData
	if !decodeData(env, &d) {
		c.sendEvent(errorEvent("消息格式不正确"))
		return
	}
	if d.ChannelID <= 0 || d.Content == "" || len(d.Content) > 5000 {
		c.sendEvent(errorEvent("消息内容不能为空且不超过 5000 字"))
		return
	}
	if !c.markSeen(d.ClientMsgID) {
		return // 幂等：重复的 client_msg_id 直接丢弃
	}
	if h.svc == nil {
		c.sendEvent(errorEvent("消息服务不可用"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	payload, err := h.svc.CreateMessage(ctx, d.ChannelID, c.userID, d.Content, d.Emotion, d.IsAnonymous)
	if err != nil {
		c.sendEvent(errorEvent("发送失败：" + err.Error()))
		return
	}
	// 落库成功后广播（含回显给发送者）。
	h.BroadcastNewMessage(d.ChannelID, payload)
}

// decodeData 解析信封的 Data 到目标结构，空 Data 视为成功（保留零值）。
func decodeData(env Envelope, dst any) bool {
	if len(env.Data) == 0 {
		return true
	}
	return json.Unmarshal(env.Data, dst) == nil
}

// errorEvent 构造一个错误信封。
func errorEvent(msg string) Envelope {
	return outbound(EventError, 0, map[string]any{"message": msg})
}
