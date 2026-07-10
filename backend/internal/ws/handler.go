// Package ws 负责 WebSocket Hub、连接管理、事件广播（经 Redis Pub/Sub）。
package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/Alike/backend/pkg/jwt"
)

// upgrader 将 HTTP 连接升级为 WebSocket。
// 生产环境应在 CheckOrigin 中校验来源；此处交由上层 CORS / Nginx 控制。
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handler 承载 WebSocket 端点的依赖。
type Handler struct {
	hub *Hub
	jwt *jwt.Manager
}

// NewHandler 创建 WebSocket handler。
func NewHandler(hub *Hub, jwtMgr *jwt.Manager) *Handler {
	return &Handler{hub: hub, jwt: jwtMgr}
}

// Handle 处理 GET /api/ws：升级连接并执行首帧鉴权。
// 鉴权约定：升级后客户端须在 5s 内发送 { "type": "auth", "data": { "token": "<JWT>" } }。
// JWT 不放 URL query，避免泄露到 Nginx 访问日志。
func (h *Handler) Handle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return // Upgrade 失败时已写入响应
	}

	userID, ok := h.authenticate(conn)
	if !ok {
		writeClose(conn, "鉴权失败")
		_ = conn.Close()
		return
	}

	client := newClient(h.hub, conn, userID)
	h.hub.register(client)
	client.sendEvent(outbound(EventAuthOK, 0, map[string]any{"user_id": userID}))

	safeGo("writePump", client.writePump)
	client.readPump()
}

// authenticate 读取首帧 auth 消息并校验 JWT，返回用户 ID。
func (h *Handler) authenticate(conn *websocket.Conn) (int64, bool) {
	_ = conn.SetReadDeadline(time.Now().Add(authTimeout))
	conn.SetReadLimit(maxMessageSize)

	_, raw, err := conn.ReadMessage()
	if err != nil {
		return 0, false
	}

	var env Envelope
	if err := json.Unmarshal(raw, &env); err != nil || env.Type != EventAuth {
		return 0, false
	}
	var d authData
	if err := json.Unmarshal(env.Data, &d); err != nil || d.Token == "" {
		return 0, false
	}
	if h.jwt == nil {
		return 0, false
	}
	claims, err := h.jwt.Parse(d.Token)
	if err != nil || claims.Type != jwt.AccessToken {
		return 0, false
	}

	// 鉴权通过后清除读超时（后续由 pongWait 心跳机制管理）。
	_ = conn.SetReadDeadline(time.Time{})
	return claims.UserID, true
}

// writeClose 向客户端发送一条错误信封后再关闭（尽力而为）。
func writeClose(conn *websocket.Conn, msg string) {
	_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
	b, _ := json.Marshal(errorEvent(msg))
	_ = conn.WriteMessage(websocket.TextMessage, b)
}
