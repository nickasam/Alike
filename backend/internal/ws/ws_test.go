package ws

import (
	"encoding/json"
	"testing"
)

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
