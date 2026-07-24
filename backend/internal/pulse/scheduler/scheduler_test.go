package scheduler

import (
	"context"
	"testing"
	"time"
)

// M0 骨架测试：Start / Stop 幂等，且 Stop 不阻塞。
func TestStartStopIdempotent(t *testing.T) {
	s := New(nil) // M0 骨架不依赖 repo，方法体不解引用

	s.Start(context.Background())
	s.Start(context.Background()) // 第二次 Start 应无操作

	done := make(chan struct{})
	go func() {
		s.Stop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Stop() 应立即返回（M0 无 worker）")
	}

	// 二次 Stop 也不应 panic
	s.Stop()
}
