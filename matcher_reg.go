package matchs


import (
	"fmt"
	"github.com/dlclark/regexp2"
)

type regRule struct {
	Tag      string
	metaData string
	reg      *regexp2.Regexp
}

func NewRegRule(str string) (*regRule, error) {
	if r, err := regexp2.Compile(str, 0); err == nil {
		return &regRule{
			Tag:      REGEXP_PREFIX + str,
			metaData: str,
			reg:      r,
		}, nil
	} else {
		return nil, fmt.Errorf("%s comiple regexp error:%s ", str, err.Error())
	}
}

func (r *regRule) MatchAll(text string) map[string]string {
	var ret = make(map[string]string, 0)
	match, _ := r.reg.FindStringMatch(text)
	//-one reg find one
	if match != nil {
		ret[r.Tag] = match.String()
		//ret = append(ret, match.String())
		//match, _ = r.reg.FindNextMatch(match) //del
	}
	//-one reg find more
	//for {
	//	if match == nil {
	//		break
	//	}
	//  ret[r.Tag] = match.String()
	//	//ret = append(ret, match.String())
	//	match, _ = r.reg.FindNextMatch(match)
	//	if match == nil {
	//		break
	//	}
	//}
	return ret
}

type RegexpMather struct {
	matchers []*regRule
}

func NewRegexpMatcher() *RegexpMather {
	return &RegexpMather{
	}
}

func (a *RegexpMather) Build(words []string) {
	for _, w := range words {
		if m, err := NewRegRule(w); err == nil {
			a.matchers = append(a.matchers, m)
		}
	}
	return
}

//Match
func (a *RegexpMather) Match(text string, onlyOne bool, repl rune) (word []string, desensitization string) {
	//desensitization = text
	for _, r := range a.matchers {
		ret := r.MatchAll(text)
		for tag, _ := range ret {
			word = append(word, tag)
		}
		//所有的正则只有命中一个就返回
		if onlyOne && len(ret) > 0 {
			return
		}
	}
	return
}
