package matchs


import "strings"

const (
	NON = iota
	AND
)

type rule struct {
	metaData string //原始数据
	words    []string
	optType  int
}

//match ####
func (r *rule) match(text string) bool {
	if r.optType == AND {
		index := -1
		//A|B|C 有大概的顺序关系，在包含多个ABC的情况下，可能不准
		for _, w := range r.words {
			last := strings.Index(text, w)
			if last > 0 && last > index {
				index = last
				continue
			} else {
				return false
			}
		}
		return true
	} else if r.optType == NON {
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

type AssembleMather struct {
	rules []*rule
}

func NewAssembleMather() *AssembleMather {
	return &AssembleMather{}
}

func (a *AssembleMather) Build(words []string) {

	for _, w := range words {

		if strings.Contains(w, "|") {
			a.rules = append(a.rules, &rule{
				metaData: w,
				words:    strings.Split(w, "|"),
				optType:  AND,
			})
		} else if strings.Contains(w, "#") {
			a.rules = append(a.rules, &rule{
				metaData: w,
				words:    strings.Split(w, "#"),
				optType:  NON,
			})
		}
	}
}

// Match onlyOne控制是否只处理一次 //不支持脱敏处理；
func (a *AssembleMather) Match(text string, onlyOne bool, repl rune) (word []string, desensitization string) {
	//desensitization = text
	for _, rule := range a.rules {
		if rule.match(text) {
			word = append(word, rule.metaData)
			if onlyOne {
				return
			}
		}
	}
	return
}
