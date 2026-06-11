package matchs

import "strings"

const (
	ruleTypeNon = iota
	ruleTypeAnd
)

type rule struct {
	raw      string
	words    []string
	ruleType int
}

// match 判断组合规则是否命中文本。
//
// ruleTypeAnd 表示 words 必须按顺序出现；ruleTypeNon 表示包含第一个词且不包含后续词。
func (r *rule) match(text string) bool {
	if r.ruleType == ruleTypeAnd {
		index := 0
		for _, w := range r.words {
			last := strings.Index(text[index:], w)
			if last < 0 {
				return false
			}
			index += last + len(w)
		}
		return true
	} else if r.ruleType == ruleTypeNon {
		if strings.Contains(text, r.words[0]) {
			for i := 1; i < len(r.words); i++ {
				if strings.Contains(text, r.words[i]) {
					return false
				}
			}
			return true
		}
	}

	return false
}

// AssembleMatcher 处理组合规则。
//
// 支持两类规则：
//   - A|B|C：A、B、C 必须按顺序出现在文本中。
//   - A#B：文本必须包含 A，且不能包含 B。
//
// AssembleMatcher 当前只返回命中的规则，不执行文本替换。
type AssembleMatcher struct {
	rules []*rule
}

// AssembleMather 是 AssembleMatcher 的兼容别名。
//
// Deprecated: 使用 AssembleMatcher。
type AssembleMather = AssembleMatcher

// NewAssembleMatcher 创建组合规则 matcher。
func NewAssembleMatcher() *AssembleMatcher {
	return &AssembleMatcher{}
}

// NewAssembleMather 创建组合规则 matcher。
//
// Deprecated: 使用 NewAssembleMatcher。
func NewAssembleMather() *AssembleMatcher {
	return NewAssembleMatcher()
}

// Build 构建组合规则列表。
func (a *AssembleMatcher) Build(words []string) {

	for _, w := range words {

		if strings.Contains(w, "|") {
			a.rules = append(a.rules, &rule{
				raw:      w,
				words:    strings.Split(w, "|"),
				ruleType: ruleTypeAnd,
			})
		} else if strings.Contains(w, "#") {
			a.rules = append(a.rules, &rule{
				raw:      w,
				words:    strings.Split(w, "#"),
				ruleType: ruleTypeNon,
			})
		}
	}
}

// Match 返回命中的组合规则。
//
// repl 参数当前不会生效；desensitization 始终返回原文。
func (a *AssembleMatcher) Match(text string, onlyOne bool, repl rune) (word []string, desensitization string) {
	desensitization = text
	for _, rule := range a.rules {
		if rule.match(text) {
			word = append(word, rule.raw)
			if onlyOne {
				return
			}
		}
	}
	return
}
