package matchs_v2

import "strings"

const (
	ruleTypeNon = iota
	ruleTypeAnd
)

type rule struct {
	raw      string
	words    []string
	ruleType int
	anchor   string
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
	rules         []*rule
	anchorIndex   map[string][]*rule
	fallbackRules []*rule
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
	if a.anchorIndex == nil {
		a.anchorIndex = make(map[string][]*rule)
	}

	for _, w := range words {
		if strings.Contains(w, "|") {
			r := &rule{
				raw:      w,
				words:    strings.Split(w, "|"),
				ruleType: ruleTypeAnd,
			}
			a.addRule(r)
		} else if strings.Contains(w, "#") {
			r := &rule{
				raw:      w,
				words:    strings.Split(w, "#"),
				ruleType: ruleTypeNon,
			}
			a.addRule(r)
		}
	}
}

func (a *AssembleMatcher) addRule(r *rule) {
	r.anchor = r.firstNonEmptyWord()
	a.rules = append(a.rules, r)
	if r.anchor == "" {
		a.fallbackRules = append(a.fallbackRules, r)
		return
	}
	a.anchorIndex[r.anchor] = append(a.anchorIndex[r.anchor], r)
}

func (r *rule) firstNonEmptyWord() string {
	for _, word := range r.words {
		if word != "" {
			return word
		}
	}
	return ""
}

func (a *AssembleMatcher) candidateRules(text string) map[*rule]struct{} {
	candidates := make(map[*rule]struct{}, len(a.fallbackRules))
	for _, rule := range a.fallbackRules {
		candidates[rule] = struct{}{}
	}
	for anchor, rules := range a.anchorIndex {
		if !strings.Contains(text, anchor) {
			continue
		}
		for _, rule := range rules {
			candidates[rule] = struct{}{}
		}
	}
	return candidates
}

// Match 返回命中的组合规则。
//
// repl 参数当前不会生效；desensitization 始终返回原文。
func (a *AssembleMatcher) Match(text string, onlyOne bool, repl rune) (word []string, desensitization string) {
	desensitization = text
	candidates := a.candidateRules(text)
	if len(candidates) == 0 {
		return
	}
	for _, rule := range a.rules {
		if _, ok := candidates[rule]; !ok {
			continue
		}
		if rule.match(text) {
			word = append(word, rule.raw)
			if onlyOne {
				return
			}
		}
	}
	return
}
