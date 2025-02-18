package base

import (
	"github.com/json-iterator/go/extra"
)

// 初始化 jsoniter
func InitJsoniter() {
	// 全局使用 snake_case
	// jsoniter 自带的 extra.LowerCaseWithUnderscores 不兼容 eventID 的情况，所以自行实现一个
	extra.SetNamingStrategy(CaseToSnake)

	// 容忍字符串和数字互转
	extra.RegisterFuzzyDecoders()
}
