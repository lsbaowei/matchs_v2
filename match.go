package matchs

import "strings"

const (
	DFA      = 0
	ASSEMBLE = 1
	REGEXP   = 2

	REGEXP_PREFIX = "reg@"
)

type MatchService struct {
	matchers map[int]Matcher
}

//初始化
func NewMatchService() *MatchService {
	return &MatchService{
		matchers: make(map[int]Matcher),
	}
}

//Build 当前支持三种配置，可以新增
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
		matcher := NewDFAMather()
		matcher.Build(dfaList)
		m.matchers[DFA] = matcher
	}

	if len(assembleList) > 0 {
		matcher := NewAssembleMather()
		matcher.Build(assembleList)
		m.matchers[ASSEMBLE] = matcher
	}

	if len(regexpList) > 0 {
		matcher := NewRegexpMatcher()
		matcher.Build(regexpList)
		m.matchers[REGEXP] = matcher
	}
}

//Match
func (m *MatchService) Match(text string, onlyOne bool, repl rune) (sensitiveWords []string, replaceText string) {
	for _, x := range m.matchers {
		ret, _ := x.Match(text, onlyOne, repl)
		for _, word := range ret {
			sensitiveWords = append(sensitiveWords, word)
		}
		//限制次数
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
