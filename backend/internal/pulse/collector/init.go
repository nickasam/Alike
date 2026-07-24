// Package collector: init 时注册所有内置 collector。
// scheduler 启动时会按 pulse_topics.collector_kind 查表调度。
package collector

func init() {
	Register(&GitHubTrending{})
	Register(&HackerNewsAI{})
}
