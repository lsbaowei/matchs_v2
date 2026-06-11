package matchs_v2

import "strings"

const (
	DFA      = 0 // DFA matcher，处理普通关键词。
	ASSEMBLE = 1 // Assemble matcher，处理 | 和 # 组合规则。
	REGEXP   = 2 // Regexp matcher，处理 reg@ 正则规则。

	REGEXP_PREFIX = "reg@"
)

var matcherOrder = []int{DFA, ASSEMBLE, REGEXP}

// MatchService 负责按固定顺序聚合不同类型的 matcher。
//
// Build 会把规则拆分为普通关键词、组合规则和正则规则；Match 会按
// DFA、ASSEMBLE、REGEXP 的顺序执行。onlyOne 为 true 时，任意 matcher
// 首次命中后立即返回。
type MatchService struct {
	matchers map[int]Matcher
}

// NewMatchService 创建一个空的匹配服务。
func NewMatchService() *MatchService {
	return &MatchService{
		matchers: make(map[int]Matcher),
	}
}

// Build 构建规则列表。
//
// 规则分类：
//   - 以 REGEXP_PREFIX 开头的规则会作为正则表达式处理。
//   - 包含 | 或 # 的规则会作为组合规则处理。
//   - 其他规则会作为普通关键词处理。
func (m *MatchService) Build(words []string) {
	var (
		dfaList      []string
		assembleList []string
		regexpList   []string
	)

	for i := 0; i < len(words); i++ {
		if strings.HasPrefix(words[i], REGEXP_PREFIX) {
			regexpList = append(regexpList, words[i][len(REGEXP_PREFIX):])
		} else if strings.Contains(words[i], "|") || strings.Contains(words[i], "#") {
			assembleList = append(assembleList, words[i])
		} else {
			dfaList = append(dfaList, words[i])
		}
	}

	if len(dfaList) > 0 {
		matcher := NewDFAMatcher()
		matcher.Build(dfaList)
		m.matchers[DFA] = matcher
	}

	if len(assembleList) > 0 {
		matcher := NewAssembleMatcher()
		matcher.Build(assembleList)
		m.matchers[ASSEMBLE] = matcher
	}

	if len(regexpList) > 0 {
		matcher := NewRegexpMatcher()
		matcher.Build(regexpList)
		m.matchers[REGEXP] = matcher
	}
}

// Match 返回命中的规则列表和替换后的文本。
//
// repl 只会用于普通关键词的 DFA 替换；组合规则和正则规则当前只返回命中的规则标识，
// 不执行脱敏替换。onlyOne 为 true 时，只返回按 matcherOrder 顺序遇到的第一类命中结果。
func (m *MatchService) Match(text string, onlyOne bool, repl rune) (sensitiveWords []string, replaceText string) {
	replaceText = text
	for _, typ := range matcherOrder {
		x, ok := m.matchers[typ]
		if !ok {
			continue
		}
		ret, replaced := x.Match(replaceText, onlyOne, repl)
		replaceText = replaced
		for _, word := range ret {
			sensitiveWords = append(sensitiveWords, word)
		}
		if onlyOne && len(ret) > 0 {
			return
		}
	}
	return
}

/*-------------other util-------------------*/

func isASCIISpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// TrimString returns s without leading and trailing ASCII space.
func TrimString(s string) string {
	for len(s) > 0 && isASCIISpace(s[0]) {
		s = s[1:]
	}
	for len(s) > 0 && isASCIISpace(s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return s
}
