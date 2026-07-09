package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

// pubSubChannel 是所有实例广播 WebSocket 事件的 Redis 频道名。
const pubSubChannel = "alike:ws:events"

// PubSub 经 Redis Pub/Sub 在多实例间广播 WebSocket 事件。
type PubSub struct {
	rdb     *redis.Client
	ctx     context.Context
	cancel  context.CancelFunc
	handler func(Envelope)
}

// NewPubSub 创建 Redis Pub/Sub 广播器。rdb 为 nil 时返回 nil（降级为纯本地广播）。
func NewPubSub(rdb *redis.Client) *PubSub {
	if rdb == nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &PubSub{rdb: rdb, ctx: ctx, cancel: cancel}
}

// OnMessage 注册从 Redis 收到跨实例事件时的本地投递回调。
func (p *PubSub) OnMessage(fn func(Envelope)) {
	p.handler = fn
}

// Start 订阅 Redis 频道并在后台循环投递事件。
func (p *PubSub) Start() {
	sub := p.rdb.Subscribe(p.ctx, pubSubChannel)
	ch := sub.Channel()
	go func() {
		defer sub.Close()
		for {
			select {
			case <-p.ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				var env Envelope
				if err := json.Unmarshal([]byte(msg.Payload), &env); err != nil {
					log.Printf("[WARN] ws: bad pubsub payload: %v", err)
					continue
				}
				if p.handler != nil {
					p.handler(env)
				}
			}
		}
	}()
}

// Publish 将事件发布到 Redis 频道，供所有实例（含本实例）消费。
func (p *PubSub) Publish(evt Envelope) error {
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	return p.rdb.Publish(p.ctx, pubSubChannel, b).Err()
}

// Close 停止订阅循环。
func (p *PubSub) Close() {
	p.cancel()
}
