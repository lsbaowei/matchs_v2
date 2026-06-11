package matchs

import (
	"fmt"
	"regexp"
)

type regRule struct {
	tag string
	reg *regexp.Regexp
}

// NewRegRule 编译一条正则规则。
//
// 传入的 str 不包含 REGEXP_PREFIX。正则语法使用 Go 标准库 regexp，
// 不支持 lookbehind、backreference 等回溯型正则特性。
func NewRegRule(str string) (*regRule, error) {
	r, err := regexp.Compile(str)
	if err != nil {
		return nil, fmt.Errorf("%s compile regexp error:%s", str, err.Error())
	}
	return &regRule{
		tag: REGEXP_PREFIX + str,
		reg: r,
	}, nil
}

// matchAll 判断正则规则是否命中文本。
//
// 当前只需要知道每条正则是否命中一次，因此返回值最多包含一个 tag。
func (r *regRule) matchAll(text string) map[string]string {
	var ret = make(map[string]string, 0)
	if r.reg.FindStringIndex(text) != nil {
		ret[r.tag] = r.reg.FindString(text)
	}
	return ret
}

// RegexpMatcher 处理 reg@ 前缀声明的正则规则。
//
// 当前只返回命中的规则 tag，不执行文本替换。
type RegexpMatcher struct {
	matchers []*regRule
}

// RegexpMather 是 RegexpMatcher 的兼容别名。
//
// Deprecated: 使用 RegexpMatcher。
type RegexpMather = RegexpMatcher

// NewRegexpMatcher 创建正则 matcher。
func NewRegexpMatcher() *RegexpMatcher {
	return &RegexpMatcher{}
}

// Build 编译正则规则列表。
//
// 无法被 Go 标准库 regexp 编译的规则会被忽略。
func (a *RegexpMatcher) Build(words []string) {
	for _, w := range words {
		if m, err := NewRegRule(w); err == nil {
			a.matchers = append(a.matchers, m)
		}
	}
}

// Match 返回命中的正则规则 tag。
//
// repl 参数当前不会生效；desensitization 始终返回原文。
func (a *RegexpMatcher) Match(text string, onlyOne bool, repl rune) (word []string, desensitization string) {
	desensitization = text
	for _, r := range a.matchers {
		ret := r.matchAll(text)
		for tag := range ret {
			word = append(word, tag)
		}
		//所有的正则只有命中一个就返回
		if onlyOne && len(ret) > 0 {
			return
		}
	}
	return
}
