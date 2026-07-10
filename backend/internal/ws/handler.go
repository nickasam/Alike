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

// Handler 承载 WebSocket 端点的依赖。
type Handler struct {
	hub      *Hub
	jwt      *jwt.Manager
	upgrader websocket.Upgrader
}

// NewHandler 创建 WebSocket handler。
// allowedOrigins 为空时放行所有来源（开发/同源部署）；非空时仅放行白名单 Origin，
// 与 REST CORS 白名单保持一致，避免跨站 WebSocket 劫持（CSWSH）。
// 无 Origin 头的非浏览器客户端一律放行（其安全性由首帧 JWT 鉴权保证）。
func NewHandler(hub *Hub, jwtMgr *jwt.Manager, allowedOrigins ...string) *Handler {
	allow := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allow[o] = struct{}{}
	}
	checkOrigin := func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" || len(allow) == 0 {
			return true
		}
		_, ok := allow[origin]
		return ok
	}
	return &Handler{
		hub: hub,
		jwt: jwtMgr,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     checkOrigin,
		},
	}
}

// Handle 处理 GET /api/ws：升级连接并执行首帧鉴权。
// 鉴权约定：升级后客户端须在 5s 内发送 { "type": "auth", "data": { "token": "<JWT>" } }。
// JWT 不放 URL query，避免泄露到 Nginx 访问日志。
func (h *Handler) Handle(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
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
