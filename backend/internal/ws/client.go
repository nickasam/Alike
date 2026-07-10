package ws

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// writeWait 是一次写操作允许的最长时间。
	writeWait = 10 * time.Second
	// pongWait 是读取下一条消息（含 pong）的最长等待时间，需大于 pingPeriod。
	pongWait = 60 * time.Second
	// pingPeriod 是服务端主动发送 ping 的间隔。
	pingPeriod = 30 * time.Second
	// authTimeout 是首帧鉴权的超时时间。
	authTimeout = 5 * time.Second
	// maxMessageSize 限制单条客户端消息大小（字节）。
	maxMessageSize = 8 * 1024
	// sendBuffer 是每个连接的出站缓冲队列长度。
	sendBuffer = 64
	// dedupCap 是每连接幂等去重集合的容量上限。
	dedupCap = 256
)

// Client 封装一个已鉴权的 WebSocket 连接。
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID int64

	mu       sync.Mutex
	channels map[int64]struct{} // 已加入（订阅）的频道
	seen     map[string]struct{} // 已处理的 client_msg_id，用于幂等去重
	closed   bool                // send 是否已关闭，防止向已关闭 channel 发送导致 panic
}

// newClient 创建一个绑定 Hub 的客户端。
func newClient(hub *Hub, conn *websocket.Conn, userID int64) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, sendBuffer),
		userID:   userID,
		channels: make(map[int64]struct{}),
		seen:     make(map[string]struct{}),
	}
}

// enqueue 将已编码的帧放入出站队列；队列满或连接已关闭则返回 false
// （调用方应关闭连接）。持有 mu 期间发送，与 closeSend 的 close 互斥，
// 避免向已关闭 channel 发送引发 panic。
func (c *Client) enqueue(payload []byte) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return false
	}
	select {
	case c.send <- payload:
		return true
	default:
		return false
	}
}

// closeSend 幂等地关闭出站队列，之后的 enqueue 均返回 false。
func (c *Client) closeSend() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	close(c.send)
}

// sendEvent 编码并投递一个服务端事件，失败静默丢弃（连接即将关闭）。
func (c *Client) sendEvent(evt Envelope) {
	b, err := json.Marshal(evt)
	if err != nil {
		return
	}
	c.enqueue(b)
}

// subscribe 记录客户端加入了某频道。
func (c *Client) subscribe(channelID int64) {
	c.mu.Lock()
	c.channels[channelID] = struct{}{}
	c.mu.Unlock()
}

// unsubscribe 移除某频道订阅。
func (c *Client) unsubscribe(channelID int64) {
	c.mu.Lock()
	delete(c.channels, channelID)
	c.mu.Unlock()
}

// subscribedChannels 返回客户端订阅的频道快照。
func (c *Client) subscribedChannels() []int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	ids := make([]int64, 0, len(c.channels))
	for id := range c.channels {
		ids = append(ids, id)
	}
	return ids
}

// isSubscribed 报告客户端是否已加入某频道。
func (c *Client) isSubscribed(channelID int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.channels[channelID]
	return ok
}

// markSeen 记录 client_msg_id；若为重复则返回 false（应跳过处理）。
// 空 id 视为不去重，恒返回 true。
func (c *Client) markSeen(id string) bool {
	if id == "" {
		return true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.seen[id]; ok {
		return false
	}
	if len(c.seen) >= dedupCap {
		// 简单清空以限制内存，牺牲极端情况下的去重精度。
		c.seen = make(map[string]struct{})
	}
	c.seen[id] = struct{}{}
	return true
}

// readPump 读取并分发客户端事件，直到连接关闭或出错。
func (c *Client) readPump() {
	defer c.hub.unregister(c)

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.hub.handleClientEvent(c, raw)
	}
}

// writePump 将出站队列写入连接，并周期性发送 ping 心跳。
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case payload, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		case <-ticker.C:
			// 应用层 ping（客户端回 {"type":"pong"}）叠加协议层 ping 控制帧，双保险。
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			b, _ := json.Marshal(Envelope{Type: EventPing})
			if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
