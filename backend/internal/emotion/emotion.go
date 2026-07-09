// Package emotion 负责情绪标签、情绪看板、趋势统计。
// 阶段一仅声明包与情绪标签常量，业务见阶段五。
package emotion

// Tag 是支持的情绪标签，与 PRD/架构文档 2.4 一致（共 8 种）。
type Tag string

const (
	TagTired   Tag = "tired"   // 😮‍💨 疲惫
	TagAngry   Tag = "angry"   // 😡 愤怒
	TagWronged Tag = "wronged" // 😢 委屈
	TagBreak   Tag = "break"   // 🤯 崩溃
	TagNumb    Tag = "numb"    // 😴 麻木
	TagQuit    Tag = "quit"    // 🔥 想润
	TagAnxious Tag = "anxious" // 😰 焦虑
	TagCheer   Tag = "cheer"   // 💪 加油
)

// AllTags 返回全部合法情绪标签。
func AllTags() []Tag {
	return []Tag{TagTired, TagAngry, TagWronged, TagBreak, TagNumb, TagQuit, TagAnxious, TagCheer}
}

// TODO 阶段五：情绪看板聚合、趋势统计。
