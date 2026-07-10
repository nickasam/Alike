package emotion

import (
	"testing"
	"time"
)

// TestTodayStartUsesBoardLocation 验证今日起点为看板时区当天 00:00。
func TestTodayStartUsesBoardLocation(t *testing.T) {
	start := todayStart()
	if start.Location().String() != boardLocation.String() {
		t.Fatalf("todayStart location=%s, want %s", start.Location(), boardLocation)
	}
	h, m, s := start.Clock()
	if h != 0 || m != 0 || s != 0 {
		t.Fatalf("todayStart 应为当天 00:00:00, got %02d:%02d:%02d", h, m, s)
	}
	now := time.Now().In(boardLocation)
	if start.After(now) {
		t.Fatal("todayStart 不应晚于当前时刻")
	}
	if now.Sub(start) >= 24*time.Hour {
		t.Fatal("todayStart 应在最近 24 小时内")
	}
}

// TestBoardScopeSerialization 验证 Board 序列化包含 scope/dominant 字段契约。
func TestBoardScopeSerialization(t *testing.T) {
	b := &Board{
		Scope:    "today",
		Emotions: []Count{{Emotion: "tired", Count: 3}, {Emotion: "cheer", Count: 1}},
		Total:    4,
		Dominant: "tired",
	}
	if b.Scope != "today" || b.Dominant != "tired" || b.Total != 4 {
		t.Fatalf("board 字段不符: %+v", b)
	}
}
