// Package ws 负责 WebSocket Hub、连接管理、事件广播（经 Redis Pub/Sub）。
// 阶段一仅提供占位 handler，Hub 实现见阶段四。
package ws

import (
	"github.com/gin-gonic/gin"

	"github.com/Alike/backend/pkg/response"
)

// Handler WebSocket 端点入口 /api/ws。
// TODO 阶段四：升级为 WebSocket，实现首帧 auth 鉴权、心跳、频道广播。
func Handler(c *gin.Context) {
	response.Success(c, gin.H{"todo": "ws", "hint": "WebSocket upgrade pending (阶段四)"})
}
