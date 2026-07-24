// Package collector: Hacker News AI 抓取实现。
// 通过 Algolia 搜索 API 并发查询多个关键词，去重后按"半衰期加权分"排序输出。
// Algolia 是 HN 官方鼓励使用的搜索接口，无需鉴权、国内可访问、稳定性高于爬取原页。
// 见 pulse-module-design.md 第 3.3 章。
package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// HackerNewsAI 抓取 HN 上关于 AI 的高赞讨论。
// 主分 = HN points；前端排序时会用半衰期时间衰减重排（新帖优先）。
type HackerNewsAI struct{}

func (h *HackerNewsAI) Kind() string { return "hackernews_ai" }

// hnAIConfig 是从 pulse_topics.collector_config 反序列化的参数。
type hnAIConfig struct {
	Keywords    []string `json:"keywords"`     // 关键词轮询列表
	MinPoints   int      `json:"min_points"`   // 分数下限（Algolia points> 过滤）
	WindowHours int      `json:"window_hours"` // 时间窗口（小时）
	Limit       int      `json:"limit"`        // 上限（默认 20）
}

// hnExtra 序列化进 pulse_items.extra，前端 NewsCard 需要 hn_id / comments / domain。
type hnExtra struct {
	HNID     string `json:"hn_id"`              // HN objectID，用于 https://news.ycombinator.com/item?id=xxx
	Comments int    `json:"comments,omitempty"` // 评论数
	Domain   string `json:"domain,omitempty"`   // 来源域名 "openai.com" / "arxiv.org"
}

// hnHit 是 Algolia response.hits[i] 里我们消费的字段子集。
type hnHit struct {
	ObjectID    string `json:"objectID"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Author      string `json:"author"`
	Points      int    `json:"points"`
	NumComments int    `json:"num_comments"`
	CreatedAt   string `json:"created_at"`
	CreatedAtI  int64  `json:"created_at_i"`
	StoryText   string `json:"story_text"`
}
type hnResp struct {
	Hits []hnHit `json:"hits"`
}

const algoliaURL = "https://hn.algolia.com/api/v1/search"

// Fetch 拉一批 HN AI 相关高赞帖。
// 策略：对每个关键词并发查一次（限最多 8 并发），按 objectID 去重后按半衰期加权分排序，取 Top N。
func (h *HackerNewsAI) Fetch(ctx context.Context, rawCfg json.RawMessage) ([]Item, error) {
	cfg := hnAIConfig{
		Keywords:    []string{"openai", "anthropic", "claude", "gpt", "llm", "gemini", "mistral", "deepseek"},
		MinPoints:   50,
		WindowHours: 72,
		Limit:       20,
	}
	if len(rawCfg) > 0 {
		_ = json.Unmarshal(rawCfg, &cfg)
	}
	if len(cfg.Keywords) == 0 {
		return nil, fmt.Errorf("hn ai: keywords empty")
	}
	if cfg.WindowHours <= 0 {
		cfg.WindowHours = 72
	}
	if cfg.MinPoints < 0 {
		cfg.MinPoints = 0
	}
	if cfg.Limit <= 0 || cfg.Limit > 100 {
		cfg.Limit = 20
	}

	since := time.Now().Add(-time.Duration(cfg.WindowHours) * time.Hour).Unix()

	// 并发查所有关键词
	type result struct {
		hits []hnHit
		err  error
	}
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	results := make([]result, len(cfg.Keywords))
	for i, kw := range cfg.Keywords {
		wg.Add(1)
		go func(idx int, keyword string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			hits, err := searchHN(ctx, keyword, cfg.MinPoints, since, cfg.Limit)
			results[idx] = result{hits: hits, err: err}
		}(i, kw)
	}
	wg.Wait()

	// 去重（同一帖子可能被多个关键词命中）+ 收集错误
	dedup := make(map[string]hnHit)
	var lastErr error
	for _, r := range results {
		if r.err != nil {
			lastErr = r.err
			continue
		}
		for _, h := range r.hits {
			if _, exists := dedup[h.ObjectID]; !exists {
				dedup[h.ObjectID] = h
			}
		}
	}
	if len(dedup) == 0 {
		if lastErr != nil {
			return nil, fmt.Errorf("hn ai: all keywords failed, last err: %w", lastErr)
		}
		return nil, fmt.Errorf("hn ai: no hits across %d keywords", len(cfg.Keywords))
	}

	// 排序：半衰期加权分 = points × exp(-ln2 × age_hours / 24)
	// 半衰期 24 小时；这样新帖比老帖优先，与主设计文档 3.5 一致。
	hits := make([]hnHit, 0, len(dedup))
	for _, h := range dedup {
		hits = append(hits, h)
	}
	now := time.Now().Unix()
	sort.Slice(hits, func(i, j int) bool {
		return hnRankedScore(hits[i], now) > hnRankedScore(hits[j], now)
	})
	if len(hits) > cfg.Limit {
		hits = hits[:cfg.Limit]
	}

	items := make([]Item, 0, len(hits))
	for _, hit := range hits {
		// URL：Algolia 有时不带 url（Ask/Show/Job HN），此时用 HN 帖子页
		outURL := hit.URL
		if outURL == "" {
			outURL = "https://news.ycombinator.com/item?id=" + hit.ObjectID
		}
		summary := hit.StoryText // Ask/Show HN 会有内容
		if summary == "" {
			// 无 story_text 时用一句副标签占位；前端 line-clamp 处理长度
			summary = ""
		}

		extra := hnExtra{
			HNID:     hit.ObjectID,
			Comments: hit.NumComments,
			Domain:   extractDomain(outURL),
		}
		extraJSON, _ := json.Marshal(extra)

		var pubAt *time.Time
		if hit.CreatedAtI > 0 {
			t := time.Unix(hit.CreatedAtI, 0)
			pubAt = &t
		}

		items = append(items, Item{
			Source:      "hackernews",
			SourceID:    hit.ObjectID,
			Title:       hit.Title,
			Summary:     strings.TrimSpace(summary),
			URL:         outURL,
			Author:      hit.Author,
			Score:       hit.Points,
			Extra:       extraJSON,
			PublishedAt: pubAt,
		})
	}
	return items, nil
}

// searchHN 单次 Algolia 请求（一个关键词一次）。
func searchHN(ctx context.Context, keyword string, minPoints int, sinceUnix int64, hitsPerPage int) ([]hnHit, error) {
	params := url.Values{}
	params.Set("query", keyword)
	params.Set("tags", "story")
	// numericFilters 用逗号分隔的复合过滤
	params.Set("numericFilters",
		"points>"+strconv.Itoa(minPoints)+",created_at_i>"+strconv.FormatInt(sinceUnix, 10))
	params.Set("hitsPerPage", strconv.Itoa(hitsPerPage))

	body, err := httpGET(ctx, algoliaURL+"?"+params.Encode(), 10*time.Second)
	if err != nil {
		return nil, err
	}
	var resp hnResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse hn json: %w", err)
	}
	return resp.Hits, nil
}

// hnRankedScore 计算某帖子的当下加权分（半衰期 24h）。
// 用于跨关键词去重后的统一排序，不入库；查询端用原始 points。
func hnRankedScore(h hnHit, nowUnix int64) float64 {
	if h.CreatedAtI <= 0 {
		return float64(h.Points)
	}
	ageHours := float64(nowUnix-h.CreatedAtI) / 3600.0
	if ageHours < 0 {
		ageHours = 0
	}
	decay := math.Exp(-math.Ln2 * ageHours / 24.0)
	return float64(h.Points) * decay
}

// extractDomain 从 URL 抽 host（去掉 www 前缀）。用于卡片显示"来源域名"。
func extractDomain(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := u.Host
	host = strings.TrimPrefix(host, "www.")
	return host
}
