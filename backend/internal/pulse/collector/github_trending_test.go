package collector

import "testing"

func TestParseIntLoose(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"5,297", 5297},
		{"3,252 stars today", 3252},
		{"128 stars this week", 128},
		{"12.7k", 12700},
		{"5.2k", 5200},
		{"1.5M", 1500000},
		{"没有数字", 0},
		{"  3,252 stars today  ", 3252},
	}
	for _, tc := range cases {
		if got := parseIntLoose(tc.in); got != tc.want {
			t.Errorf("parseIntLoose(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestExtractHexColor(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"background-color: #dea584", "#dea584"},
		{"background-color: #f1e05a;", "#f1e05a"},
		{"background-color:#00add8;", "#00add8"},
		{"color: red", ""},
		{"", ""},
		{"background-color: #ab", ""}, // 太短
	}
	for _, tc := range cases {
		if got := extractHexColor(tc.in); got != tc.want {
			t.Errorf("extractHexColor(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSplitOwnerRepo(t *testing.T) {
	cases := []struct {
		in, wantOwner, wantRepo string
	}{
		{"microsoft/vscode", "microsoft", "vscode"},
		{"anthropics/claude-code", "anthropics", "claude-code"},
		{"single", "", ""},
		{"a/b/c", "a", "b"}, // 多余段丢弃
	}
	for _, tc := range cases {
		o, r := splitOwnerRepo(tc.in)
		if o != tc.wantOwner || r != tc.wantRepo {
			t.Errorf("splitOwnerRepo(%q) = (%q,%q), want (%q,%q)", tc.in, o, r, tc.wantOwner, tc.wantRepo)
		}
	}
}

// 确认 init 注册成功。
func TestGitHubTrendingRegistered(t *testing.T) {
	c, ok := Get("github_trending")
	if !ok {
		t.Fatal("github_trending not registered by init()")
	}
	if c.Kind() != "github_trending" {
		t.Fatalf("Kind() = %q", c.Kind())
	}
}
