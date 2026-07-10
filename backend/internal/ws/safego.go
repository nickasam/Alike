package ws

import "log"

// safeGo 在独立 goroutine 中运行 fn，并用 recover 兜底 panic。
// WebSocket 的 writePump、Redis 订阅循环等都脱离 Gin 的 Recovery 中间件，
// 任一未捕获的 panic（坏 payload、nil 解引用）会击穿整个进程，拖垮所有连接与实例。
// 统一经此 helper 启动，保证单个连接/事件的 panic 只记录日志、不影响进程。
func safeGo(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERROR] ws: goroutine %q panic recovered: %v", name, r)
			}
		}()
		fn()
	}()
}

// safeInvoke 同步执行 fn 并用 recover 兜底，用于在已处于 goroutine 中的
// 逐事件回调（如 Redis 订阅循环里的单条投递），使一条坏事件不会终止整个循环。
func safeInvoke(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] ws: %q panic recovered: %v", name, r)
		}
	}()
	fn()
}
