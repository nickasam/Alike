package collector

import "testing"

func TestExtractDomain(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"https://openai.com/blog/introducing-gpt-5", "openai.com"},
		{"https://www.techcrunch.com/2026/07/23/ai-news", "techcrunch.com"},
		{"https://news.ycombinator.com/item?id=123", "news.ycombinator.com"},
		{"https://arxiv.org/abs/2401.12345", "arxiv.org"},
		{"", ""},
		{"not a url", ""},
	}
	for _, tc := range cases {
		if got := extractDomain(tc.in); got != tc.want {
			t.Errorf("extractDomain(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestHNRankedScore(t *testing.T) {
	now := int64(1_784_800_000) // 参考时点

	// 全新帖子（0 小时前）应等于 points × 1
	fresh := hnHit{Points: 100, CreatedAtI: now}
	if got := hnRankedScore(fresh, now); got < 99.9 || got > 100.1 {
		t.Errorf("fresh score = %f, want ≈100", got)
	}

	// 24 小时前应等于 points × 0.5
	oneDay := hnHit{Points: 100, CreatedAtI: now - 24*3600}
	if got := hnRankedScore(oneDay, now); got < 49 || got > 51 {
		t.Errorf("1-day-old score = %f, want ≈50", got)
	}

	// 48 小时前应等于 points × 0.25
	twoDay := hnHit{Points: 100, CreatedAtI: now - 48*3600}
	if got := hnRankedScore(twoDay, now); got < 24 || got > 26 {
		t.Errorf("2-day-old score = %f, want ≈25", got)
	}

	// CreatedAtI=0 时降级到原始分
	unknownTime := hnHit{Points: 42, CreatedAtI: 0}
	if got := hnRankedScore(unknownTime, now); got != 42 {
		t.Errorf("unknown-time score = %f, want 42", got)
	}
}

func TestHackerNewsAIRegistered(t *testing.T) {
	c, ok := Get("hackernews_ai")
	if !ok {
		t.Fatal("hackernews_ai not registered by init()")
	}
	if c.Kind() != "hackernews_ai" {
		t.Fatalf("Kind() = %q", c.Kind())
	}
}
