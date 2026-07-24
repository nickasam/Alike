// Package scheduler 是 pulse 的定时抓取调度器。
// M1：为每个 is_active=true 的专题起独立 worker goroutine，按 refresh_interval_min 周期性 tick。
// 失败不覆写数据库；成功后 upsert Item 并更新 last_fetched_at。
package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Alike/backend/internal/pulse"
	"github.com/Alike/backend/internal/pulse/collector"
)

// Scheduler 按每个专题的 refresh_interval_min 周期性触发对应 collector。
// 单实例部署下用 time.Ticker + goroutine（不引入 cron 依赖）。
// 多实例部署时需外接分布式锁（见 pulse-module-design.md 附录 B.3），当前阶段不做。
type Scheduler struct {
	repo *pulse.Repository

	mu     sync.Mutex
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New 创建调度器实例。集成到应用启动时调用。
func New(repo *pulse.Repository) *Scheduler {
	return &Scheduler{repo: repo}
}

// Start 启动调度循环：
//  1. 读 pulse_topics（is_active=true）
//  2. 每个 topic 起一个 worker goroutine，先立即跑一次不等 tick，然后按 refresh_interval_min 循环
//  3. 起一个 cleanup goroutine 每 24h 删除 7 天前的老 items
//  4. 幂等：重复 Start 只启动一次
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancel != nil {
		return // 已启动，幂等
	}
	sub, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	if s.repo == nil {
		log.Printf("pulse.scheduler: repo=nil, skipping worker startup (unit-test mode)")
		return
	}

	topics, err := s.repo.ListActiveTopics(sub)
	if err != nil {
		log.Printf("pulse.scheduler: list topics failed: %v", err)
		return
	}
	started := 0
	for _, t := range topics {
		c, ok := collector.Get(t.CollectorKind)
		if !ok {
			log.Printf("pulse.scheduler: topic=%q collector_kind=%q not registered, skip", t.Slug, t.CollectorKind)
			continue
		}
		s.wg.Add(1)
		go s.runTopic(sub, t, c)
		started++
	}
	log.Printf("pulse.scheduler: started %d workers (out of %d active topics)", started, len(topics))

	// 清理任务：每 24 小时删除 7 天前的老 items（外部数据镜像，不做软删除）。
	s.wg.Add(1)
	go s.runCleanup(sub)
}

// Stop 优雅停止：取消所有 worker 上下文，等待退出（最多 5 秒兜底）。
func (s *Scheduler) Stop() {
	s.mu.Lock()
	cancel := s.cancel
	s.cancel = nil
	s.mu.Unlock()
	if cancel == nil {
		return
	}
	cancel()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		log.Printf("pulse.scheduler: shutdown timeout, some workers may still be running")
	}
	log.Printf("pulse.scheduler: stopped")
}

// runTopic 是单个专题的 worker：启动立即跑一次，然后按周期 tick。
func (s *Scheduler) runTopic(ctx context.Context, t *pulse.Topic, c collector.Collector) {
	defer s.wg.Done()

	// 启动即跑一次，不等第一个 tick（M0 阶段用户看空页太久）
	s.fetchOnce(ctx, t, c)

	interval := time.Duration(t.RefreshIntervalMin) * time.Minute
	if interval < time.Minute {
		interval = 60 * time.Minute // 兜底
	}
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			s.fetchOnce(ctx, t, c)
		}
	}
}

// fetchOnce 触发一次 collector.Fetch，成功则 upsert 每条 Item + 标记 last_fetched_at；
// 失败则只记 last_error，不覆写既有 pulse_items 数据（用户仍看上次成功的快照）。
func (s *Scheduler) fetchOnce(ctx context.Context, t *pulse.Topic, c collector.Collector) {
	start := time.Now()
	items, err := c.Fetch(ctx, t.CollectorConfig)
	if err != nil {
		log.Printf("pulse.scheduler: topic=%q fetch failed: %v", t.Slug, err)
		if e := s.repo.MarkTopicError(ctx, t.ID, err.Error()); e != nil {
			log.Printf("pulse.scheduler: topic=%q mark error failed: %v", t.Slug, e)
		}
		return
	}

	upserted := 0
	for _, it := range items {
		row := &pulse.Item{
			TopicID:     t.ID,
			Source:      it.Source,
			SourceID:    it.SourceID,
			Title:       it.Title,
			Summary:     it.Summary,
			URL:         it.URL,
			Author:      it.Author,
			Score:       it.Score,
			Extra:       it.Extra,
			PublishedAt: it.PublishedAt,
		}
		if err := s.repo.UpsertItem(ctx, row); err != nil {
			log.Printf("pulse.scheduler: topic=%q upsert %q failed: %v", t.Slug, it.SourceID, err)
			continue
		}
		upserted++
	}
	if err := s.repo.MarkTopicFetched(ctx, t.ID); err != nil {
		log.Printf("pulse.scheduler: topic=%q mark fetched failed: %v", t.Slug, err)
	}
	log.Printf("pulse.scheduler: topic=%q fetched=%d upserted=%d took=%s",
		t.Slug, len(items), upserted, time.Since(start).Round(time.Millisecond))
}

// cleanupWindow 是 pulse_items 的保留窗口。7 天前的数据由 runCleanup 定期物理删除。
// 参见 pulse-module-design.md 第 3.6 节。
const cleanupWindow = 7 * 24 * time.Hour

// cleanupInterval 是 runCleanup 的循环间隔。每天跑一次即可（数据量小，无需更密）。
const cleanupInterval = 24 * time.Hour

// runCleanup 定时删除 pulse_items 中超过 cleanupWindow 的老数据。
// 启动即跑一次不等 tick（新部署时可能有历史留存），随后每 24 小时一次。
func (s *Scheduler) runCleanup(ctx context.Context) {
	defer s.wg.Done()

	// 立即跑一次
	s.cleanupOnce(ctx)

	tick := time.NewTicker(cleanupInterval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			s.cleanupOnce(ctx)
		}
	}
}

// cleanupOnce 触发一次 pulse_items 清理，删除 cleanupWindow 之前的行。
func (s *Scheduler) cleanupOnce(ctx context.Context) {
	deleted, err := s.repo.CleanupOldItems(ctx, cleanupWindow)
	if err != nil {
		log.Printf("pulse.scheduler: cleanup failed: %v", err)
		return
	}
	if deleted > 0 {
		log.Printf("pulse.scheduler: cleanup removed %d rows older than %s", deleted, cleanupWindow)
	}
}
