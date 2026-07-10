package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// TestClientEnqueueAfterCloseNoPanic 回归测试：unregister 关闭 send channel 后，
// 并发的广播 enqueue 不得 panic（曾因 close 与 send 竞态导致 send on closed channel）。
func TestClientEnqueueAfterCloseNoPanic(t *testing.T) {
	c := newClient(nil, nil, 1)
	c.closeSend()
	if c.enqueue([]byte("x")) {
		t.Fatal("enqueue after closeSend should return false, not send on closed channel")
	}
	// 二次关闭需幂等，不 panic。
	c.closeSend()
}

// TestClientConcurrentEnqueueAndClose 在 -race 下验证并发 enqueue/closeSend 无数据竞争与 panic。
func TestClientConcurrentEnqueueAndClose(t *testing.T) {
	for i := 0; i < 50; i++ {
		c := newClient(nil, nil, 1)
		// 排空 send 的 goroutine，避免缓冲满掩盖竞态。
		done := make(chan struct{})
		go func() {
			for range c.send {
			}
			close(done)
		}()

		var wg sync.WaitGroup
		wg.Add(11)
		for j := 0; j < 10; j++ {
			go func() { defer wg.Done(); c.enqueue([]byte("frame")) }()
		}
		go func() { defer wg.Done(); c.closeSend() }()
		wg.Wait()
		<-done
	}
}

func TestClientMarkSeenDedup(t *testing.T) {
	c := newClient(nil, nil, 1)
	if !c.markSeen("uuid-1") {
		t.Fatal("first occurrence should be accepted")
	}
	if c.markSeen("uuid-1") {
		t.Fatal("duplicate client_msg_id should be rejected")
	}
	// 空 id 不去重，恒接受。
	if !c.markSeen("") || !c.markSeen("") {
		t.Fatal("empty client_msg_id must never dedup")
	}
}

func TestClientMarkSeenEviction(t *testing.T) {
	c := newClient(nil, nil, 1)
	for i := 0; i < dedupCap+5; i++ {
		c.markSeen(string(rune('a')) + itoaTest(i))
	}
	if len(c.seen) > dedupCap {
		t.Fatalf("seen set grew beyond cap: %d > %d", len(c.seen), dedupCap)
	}
}

func TestClientSubscribe(t *testing.T) {
	c := newClient(nil, nil, 1)
	c.subscribe(10)
	c.subscribe(20)
	if !c.isSubscribed(10) || !c.isSubscribed(20) {
		t.Fatal("subscribed channels not tracked")
	}
	c.unsubscribe(10)
	if c.isSubscribed(10) {
		t.Fatal("channel 10 should be unsubscribed")
	}
	if got := c.subscribedChannels(); len(got) != 1 {
		t.Fatalf("subscribedChannels len=%d, want 1", len(got))
	}
}

func TestOutboundEnvelope(t *testing.T) {
	env := outbound(EventNewMessage, 42, map[string]any{"id": 7})
	if env.Type != EventNewMessage || env.ChannelID != 42 {
		t.Fatalf("bad envelope: %+v", env)
	}
	var decoded map[string]any
	if err := json.Unmarshal(env.Data, &decoded); err != nil {
		t.Fatalf("data not valid json: %v", err)
	}
	if decoded["id"].(float64) != 7 {
		t.Errorf("data id=%v, want 7", decoded["id"])
	}
}

func TestDecodeData(t *testing.T) {
	env := Envelope{Data: json.RawMessage(`{"channel_id":5}`)}
	var d channelData
	if !decodeData(env, &d) || d.ChannelID != 5 {
		t.Fatalf("decode failed: %+v", d)
	}
	// 空 Data 视为成功（零值）。
	if !decodeData(Envelope{}, &d) {
		t.Fatal("empty data should decode as success")
	}
	// 非法 JSON 返回 false。
	if decodeData(Envelope{Data: json.RawMessage(`{`)}, &d) {
		t.Fatal("malformed data should fail")
	}
}

func TestHubJoinLeaveChannel(t *testing.T) {
	h := NewHub(nil, nil)
	c := newClient(h, nil, 1)
	h.register(c)
	h.joinChannel(c, 100)

	h.mu.RLock()
	_, inChannel := h.channels[100][c]
	h.mu.RUnlock()
	if !inChannel {
		t.Fatal("client not in channel set after join")
	}

	h.leaveChannel(c, 100)
	h.mu.RLock()
	set := h.channels[100]
	h.mu.RUnlock()
	if set != nil {
		t.Fatal("empty channel set should be cleaned up")
	}
}

func TestHubUnregisterCleansChannels(t *testing.T) {
	h := NewHub(nil, nil)
	c := newClient(h, nil, 1)
	h.register(c)
	h.joinChannel(c, 1)
	h.joinChannel(c, 2)
	h.unregister(c)

	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.clients) != 0 || len(h.channels) != 0 {
		t.Fatalf("unregister left residue: clients=%d channels=%d", len(h.clients), len(h.channels))
	}
}

// itoaTest 是测试内联的整型转字符串工具（避免 import strconv 触发未用告警混淆）。
func itoaTest(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// TestSafeGoRecoversPanic 验证 safeGo 内 panic 被兜底，不会传播导致进程崩溃。
func TestSafeGoRecoversPanic(t *testing.T) {
	done := make(chan struct{})
	safeGo("test-panic", func() {
		defer close(done)
		panic("boom")
	})
	select {
	case <-done:
		// 未 panic 到测试进程即通过（recover 生效）。
	case <-time.After(2 * time.Second):
		t.Fatal("safeGo goroutine 未完成")
	}
}

// TestSafeInvokeRecoversPanic 验证 safeInvoke 同步兜底 panic。
func TestSafeInvokeRecoversPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("safeInvoke 未兜底 panic: %v", r)
		}
	}()
	safeInvoke("test", func() { panic("boom") })
}

// TestHubShutdownClosesClients 验证 Shutdown 关闭所有在线客户端出站队列（幂等）。
func TestHubShutdownClosesClients(t *testing.T) {
	h := NewHub(nil, nil)
	c := newClient(h, nil, 1)
	h.register(c)

	h.Shutdown()
	// 关闭后 enqueue 应返回 false（send 已关闭）。
	if c.enqueue([]byte("x")) {
		t.Fatal("Shutdown 后 enqueue 应失败")
	}
	// 再次 Shutdown 应幂等，不 panic。
	h.Shutdown()
}

// TestBroadcastEmpathyDeliversToChannel 验证共情广播投递给订阅频道的客户端，
// 且 payload 含 message_id 与 empathy_count（与前端 applyEmpathy 对齐）。
func TestBroadcastEmpathyDeliversToChannel(t *testing.T) {
	h := NewHub(nil, nil) // pubsub 为 nil → 走本地广播
	c := newClient(h, nil, 1)
	h.register(c)
	h.joinChannel(c, 100)

	h.BroadcastEmpathy(100, 7, 5)

	select {
	case raw := <-c.send:
		var env Envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if env.Type != EventEmpathy || env.ChannelID != 100 {
			t.Fatalf("type=%s channel=%d, want empathy/100", env.Type, env.ChannelID)
		}
		var d map[string]any
		_ = json.Unmarshal(env.Data, &d)
		if d["message_id"].(float64) != 7 || d["empathy_count"].(float64) != 5 {
			t.Fatalf("payload=%v, want message_id=7 empathy_count=5", d)
		}
	default:
		t.Fatal("订阅客户端未收到 empathy 广播")
	}
}

// TestBroadcastMessageDeletedDelivers 验证删除广播投递且 payload 含 message_id。
func TestBroadcastMessageDeletedDelivers(t *testing.T) {
	h := NewHub(nil, nil)
	c := newClient(h, nil, 1)
	h.register(c)
	h.joinChannel(c, 200)

	h.BroadcastMessageDeleted(200, 42)

	select {
	case raw := <-c.send:
		var env Envelope
		_ = json.Unmarshal(raw, &env)
		if env.Type != EventMessageDeleted {
			t.Fatalf("type=%s, want message_deleted", env.Type)
		}
		var d map[string]any
		_ = json.Unmarshal(env.Data, &d)
		if d["message_id"].(float64) != 42 {
			t.Fatalf("payload=%v, want message_id=42", d)
		}
	default:
		t.Fatal("订阅客户端未收到 message_deleted 广播")
	}
}

// TestCheckOriginWhitelist 验证 WS 升级的来源校验：空白名单放行所有，
// 非空白名单仅放行命中项，无 Origin 头（非浏览器客户端）一律放行。
func TestCheckOriginWhitelist(t *testing.T) {
	mkReq := func(origin string) *http.Request {
		r := httptest.NewRequest(http.MethodGet, "/api/ws", nil)
		if origin != "" {
			r.Header.Set("Origin", origin)
		}
		return r
	}

	// 空白名单：全部放行。
	hOpen := NewHandler(nil, nil)
	if !hOpen.upgrader.CheckOrigin(mkReq("https://evil.example")) {
		t.Fatal("空白名单应放行任意 Origin")
	}

	// 白名单模式。
	h := NewHandler(nil, nil, "https://alike.example")
	if !h.upgrader.CheckOrigin(mkReq("https://alike.example")) {
		t.Fatal("白名单 Origin 应放行")
	}
	if h.upgrader.CheckOrigin(mkReq("https://evil.example")) {
		t.Fatal("非白名单 Origin 应拒绝")
	}
	if !h.upgrader.CheckOrigin(mkReq("")) {
		t.Fatal("无 Origin 头（非浏览器客户端）应放行")
	}
}
