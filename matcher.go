package matchs_v2

// Matcher 定义所有匹配器的统一接口。
//
// Build 负责构建内部匹配结构；Match 负责返回命中的规则列表和处理后的文本。
// 不同 matcher 对 repl 的支持不同：当前只有 DFAMatcher 会执行文本替换。
type Matcher interface {
	// Build 构建 matcher 内部规则。
	Build(words []string)

	// Match 返回命中的规则列表和处理后的文本。
	Match(text string, onlyOne bool, repl rune) ([]string, string)
}
