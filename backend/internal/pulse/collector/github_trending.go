// Package collector: GitHub Trending 抓取实现。
// 双策略：
//   1. 优先抓 https://github.com/trending?since=daily 的 HTML（goquery 解析，最精准）
//   2. 失败时降级到 https://api.github.com/search/repositories （国内可访问、稳定但不完全等价 Trending）
// GitHub 官方没有 Trending API，Trending 页 HTML 结构相对稳定；如有变更 fallback 兜底。
package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GitHubTrending 抓取 GitHub 每日 Star 榜。
// 主分 = 今日 +N stars；副字段（语言、总 star、fork 数）落 Extra JSONB。
type GitHubTrending struct{}

func (g *GitHubTrending) Kind() string { return "github_trending" }

// ghTrendingConfig 是从 pulse_topics.collector_config 反序列化的参数。
type ghTrendingConfig struct {
	Since string `json:"since"` // daily / weekly / monthly（默认 daily）
	Limit int    `json:"limit"` // 上限 25（默认 25）
}

// githubExtra 是仓库卡片的次要字段，序列化进 pulse_items.extra (JSONB)。
type githubExtra struct {
	Lang       string `json:"lang,omitempty"`       // "Go" / "TypeScript"
	LangColor  string `json:"lang_color,omitempty"` // "#00ADD8"
	TotalStars int    `json:"total_stars,omitempty"`
	Forks      int    `json:"forks,omitempty"`
}

const trendingURL = "https://github.com/trending"

// 伪装成常规浏览器：GitHub 会挡明显的爬虫 UA。
const chromeUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 " +
	"(KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"

// Fetch 拉取一次 GitHub Trending。
// 策略：先试 trending 页 HTML；失败（网络墙 / 结构变更）则降级 api.github.com/search。
// 失败时返回 error；scheduler 收到 error 不会覆写数据库。
func (g *GitHubTrending) Fetch(ctx context.Context, rawCfg json.RawMessage) ([]Item, error) {
	cfg := ghTrendingConfig{Since: "daily", Limit: 25}
	if len(rawCfg) > 0 {
		_ = json.Unmarshal(rawCfg, &cfg) // 反序列化失败时用默认值
	}
	if cfg.Limit <= 0 || cfg.Limit > 100 {
		cfg.Limit = 25
	}

	// 主策略：Trending HTML（最精准）
	items, errHTML := fetchTrendingHTML(ctx, cfg)
	if errHTML == nil && len(items) > 0 {
		return items, nil
	}

	// 降级：api.github.com/search（国内可达）
	items, errAPI := fetchTrendingViaAPI(ctx, cfg)
	if errAPI == nil && len(items) > 0 {
		return items, nil
	}

	// 两条路都失败
	return nil, fmt.Errorf("github trending: html=%v api=%v", errHTML, errAPI)
}

// fetchTrendingHTML 抓 https://github.com/trending 页面。
func fetchTrendingHTML(ctx context.Context, cfg ghTrendingConfig) ([]Item, error) {
	// 3 次重试 + 指数退避（500ms/1s/2s）
	var (
		body []byte
		err  error
	)
	backoff := 500 * time.Millisecond
	for attempt := 0; attempt < 3; attempt++ {
		body, err = httpGET(ctx, trendingURL+"?since="+cfg.Since, 10*time.Second)
		if err == nil {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
			backoff *= 2
		}
	}
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	items := make([]Item, 0, cfg.Limit)
	doc.Find("article.Box-row").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if len(items) >= cfg.Limit {
			return false
		}
		item, ok := parseTrendingRow(s)
		if !ok {
			return true
		}
		items = append(items, item)
		return true
	})
	if len(items) == 0 {
		return nil, fmt.Errorf("no rows parsed (HTML structure may have changed)")
	}
	return items, nil
}

// searchRepo 是 api.github.com/search/repositories 响应里 items[i] 的字段子集。
type searchRepo struct {
	FullName        string `json:"full_name"`
	HTMLURL         string `json:"html_url"`
	Description     string `json:"description"`
	Language        string `json:"language"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	Owner           struct {
		Login string `json:"login"`
	} `json:"owner"`
}
type searchResp struct {
	Items []searchRepo `json:"items"`
}

// fetchTrendingViaAPI 从 api.github.com/search/repositories 拉取"最近创建 + star 高"仓库。
// 不完全等价 Trending（Trending 用非公开算法计算日增 star），但国内可稳定访问且能给出有意义的 top 列表。
// 用作 HTML 抓取失败的降级方案。
func fetchTrendingViaAPI(ctx context.Context, cfg ghTrendingConfig) ([]Item, error) {
	// 近 30 天创建 + 按 star 排序：能覆盖大多数"当红"的新仓库
	// 若拉不到足够条目再退到 pushed:>近 7 天的老仓库
	primary := "created:>" + timeAgo(30*24*time.Hour) + " sort:stars"
	url := "https://api.github.com/search/repositories?q=" +
		strings.ReplaceAll(primary, " ", "+") +
		"&order=desc&per_page=" + strconv.Itoa(cfg.Limit)

	body, err := httpGET(ctx, url, 10*time.Second)
	if err != nil {
		return nil, err
	}
	var resp searchResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse api json: %w", err)
	}
	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("api returned 0 items")
	}

	items := make([]Item, 0, len(resp.Items))
	for _, r := range resp.Items {
		parts := strings.SplitN(r.FullName, "/", 2)
		if len(parts) != 2 {
			continue
		}
		extra := githubExtra{
			Lang:       r.Language,
			LangColor:  languageColor(r.Language),
			TotalStars: r.StargazersCount,
			Forks:      r.ForksCount,
		}
		extraJSON, _ := json.Marshal(extra)
		items = append(items, Item{
			Source:   "github",
			SourceID: r.FullName,
			Title:    r.FullName,
			Summary:  r.Description,
			URL:      r.HTMLURL,
			Author:   r.Owner.Login,
			// Search API 拿不到"今日 +N stars"精确值，
			// 退而用总 star 数作为主分（同 Trending 排序方向：数字越大越靠前）
			Score: r.StargazersCount,
			Extra: extraJSON,
		})
	}
	return items, nil
}

// timeAgo 生成 YYYY-MM-DD 格式的相对日期，供 GitHub search 查询语法使用。
func timeAgo(d time.Duration) string {
	return time.Now().UTC().Add(-d).Format("2006-01-02")
}

// languageColor 给常见语言一个固定色（GitHub 官方 linguist 色板节选）。
// api.github.com/search 不返回语言色，前端需要色点渲染，此处兜底。
func languageColor(lang string) string {
	switch lang {
	case "TypeScript":
		return "#3178c6"
	case "JavaScript":
		return "#f1e05a"
	case "Python":
		return "#3572a5"
	case "Go":
		return "#00add8"
	case "Rust":
		return "#dea584"
	case "Java":
		return "#b07219"
	case "C++":
		return "#f34b7d"
	case "C":
		return "#555555"
	case "C#":
		return "#178600"
	case "Swift":
		return "#f05138"
	case "Ruby":
		return "#701516"
	case "PHP":
		return "#4f5d95"
	case "Kotlin":
		return "#a97bff"
	case "Shell":
		return "#89e051"
	case "HTML":
		return "#e34c26"
	case "CSS":
		return "#563d7c"
	case "Vue":
		return "#41b883"
	case "Dart":
		return "#00b4ab"
	case "Zig":
		return "#ec915c"
	case "Lua":
		return "#000080"
	}
	return ""
}

// parseTrendingRow 从单个 <article class="Box-row"> 中抽字段。
// 失败时返回 ok=false（不 panic，跳过异常行）。
func parseTrendingRow(s *goquery.Selection) (Item, bool) {
	// 仓库名：h2 > a 的 href 是 /owner/repo
	link := s.Find("h2 a").First()
	href, exists := link.Attr("href")
	if !exists || len(href) < 2 {
		return Item{}, false
	}
	owner, repo := splitOwnerRepo(strings.TrimPrefix(href, "/"))
	if owner == "" || repo == "" {
		return Item{}, false
	}
	sourceID := owner + "/" + repo

	// 描述：p.color-fg-muted.my-1（可能不存在）
	summary := strings.TrimSpace(s.Find("p.color-fg-muted").First().Text())

	// 语言 + 语言色
	extra := githubExtra{}
	langNode := s.Find(`span[itemprop="programmingLanguage"]`).First()
	extra.Lang = strings.TrimSpace(langNode.Text())
	if style, ok := s.Find("span.repo-language-color").First().Attr("style"); ok {
		extra.LangColor = extractHexColor(style)
	}

	// 总 star 数：a[href$="/stargazers"]
	extra.TotalStars = parseIntLoose(s.Find(`a[href$="/stargazers"]`).First().Text())
	// fork 数：a[href$="/forks"]
	extra.Forks = parseIntLoose(s.Find(`a[href$="/forks"]`).First().Text())

	// 今日 +N stars：span.d-inline-block.float-sm-right 里含 "3,252 stars today"（可能是 weekly/monthly 时文案变化）
	todayText := strings.TrimSpace(s.Find("span.d-inline-block.float-sm-right").First().Text())
	todayStars := parseIntLoose(todayText)

	// Extra 序列化
	extraJSON, err := json.Marshal(extra)
	if err != nil {
		extraJSON = []byte("{}")
	}

	return Item{
		Source:   "github",
		SourceID: sourceID,
		Title:    sourceID, // 前端渲染 owner/repo
		Summary:  summary,
		URL:      "https://github.com/" + sourceID,
		Author:   owner,
		Score:    todayStars,
		Extra:    extraJSON,
	}, true
}

// splitOwnerRepo 把 "owner/repo" 拆开。多余段一律丢。
func splitOwnerRepo(path string) (string, string) {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

// extractHexColor 从 `background-color: #dea584` 之类的 style 值里抽 hex 颜色。
func extractHexColor(style string) string {
	idx := strings.Index(style, "#")
	if idx < 0 {
		return ""
	}
	end := idx + 1
	for end < len(style) {
		c := style[end]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
			end++
		} else {
			break
		}
	}
	if end-idx < 4 { // "#" + 至少 3 位
		return ""
	}
	return style[idx:end]
}

// parseIntLoose 从 "5,297" / "3,252 stars today" / "5.2k" 之类的文本里抽整数。
// 遇到 k/K → *1000，m/M → *1_000_000（保留一位小数精度）。
func parseIntLoose(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	// 取第一段数字（可能含小数点和逗号）
	var buf strings.Builder
	sawDigit := false
	hasDot := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			buf.WriteByte(c)
			sawDigit = true
			continue
		}
		if c == ',' && sawDigit {
			continue // 千分位分隔符
		}
		if c == '.' && sawDigit && !hasDot {
			buf.WriteByte(c)
			hasDot = true
			continue
		}
		if sawDigit {
			// 数字段结束，检查后缀
			rest := strings.ToLower(strings.TrimSpace(s[i:]))
			multiplier := 1
			if strings.HasPrefix(rest, "k") {
				multiplier = 1000
			} else if strings.HasPrefix(rest, "m") {
				multiplier = 1000000
			}
			return applyMultiplier(buf.String(), multiplier)
		}
	}
	if buf.Len() == 0 {
		return 0
	}
	return applyMultiplier(buf.String(), 1)
}

func applyMultiplier(numStr string, mult int) int {
	if strings.Contains(numStr, ".") {
		f, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0
		}
		return int(f * float64(mult))
	}
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}
	return n * mult
}

// httpGET 抓一次 HTML/JSON；带 UA 伪装 + 可控超时。
func httpGET(ctx context.Context, url string, timeout time.Duration) ([]byte, error) {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	sub, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(sub, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", chromeUA)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,application/json;q=0.8,*/*;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
