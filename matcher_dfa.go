package matchs

// DFAMatcher 使用 trie 结构匹配普通关键词。
//
// 当前实现会返回所有扫描到的关键词；当 onlyOne 为 true 时，首次命中后立即返回。
// 替换逻辑按 rune 下标执行，适用于中文等多字节字符。
type DFAMatcher struct {
	root *DFANode
}

// NewDFAMatcher 创建普通关键词 matcher。
func NewDFAMatcher() *DFAMatcher {
	return &DFAMatcher{
		root: &DFANode{
			End: false,
		},
	}
}

// NewDFAMather 创建普通关键词 matcher。
//
// Deprecated: 使用 NewDFAMatcher。
func NewDFAMather() *DFAMatcher {
	return NewDFAMatcher()
}

// Build 构造 DFA trie。
//
// 空字符串关键词会被忽略，避免根节点被误标记为命中节点。
func (d *DFAMatcher) Build(words []string) {
	for _, item := range words {
		d.root.AddWord(item)
	}
}

// Match 查找普通关键词并返回替换后的文本。
//
// repl 会替换命中的 rune 区间。存在重叠关键词时，当前实现会按扫描顺序返回多个命中结果。
func (d *DFAMatcher) Match(text string, onlyOne bool, repl rune) (sensitiveWords []string, replaceText string) {
	if d.root == nil {
		replaceText = text
		return
	}

	textChars := []rune(text)
	textCharsCopy := make([]rune, len(textChars))
	copy(textCharsCopy, textChars)

	length := len(textChars)
	for i := 0; i < length; i++ {
		// root 本身没有 key，下一层节点才表示关键词的第一个字符。
		temp := d.root.FindChild(textChars[i])
		if temp == nil {
			continue
		}
		j := i + 1
		for ; j < length && temp != nil; j++ {
			if temp.End {
				sensitiveWords = append(sensitiveWords, string(textChars[i:j]))
				replaceRune(textCharsCopy, repl, i, j)
				if onlyOne {
					replaceText = string(textCharsCopy)
					return
				}
			}
			temp = temp.FindChild(textChars[j])
		}

		if j == length && temp != nil && temp.End {
			sensitiveWords = append(sensitiveWords, string(textChars[i:length]))
			replaceRune(textCharsCopy, repl, i, length)
			if onlyOne {
				replaceText = string(textCharsCopy)
				return
			}
		}
	}
	replaceText = string(textCharsCopy)
	return
}

// replaceRune 将 [begin, end) 范围内的 rune 替换为 replaceChar。
func replaceRune(chars []rune, replaceChar rune, begin int, end int) {
	for i := begin; i < end; i++ {
		chars[i] = replaceChar
	}
}

// DFANode 是 DFA trie 的节点。
type DFANode struct {
	End  bool
	Next map[rune]*DFANode
}

// AddWord 把一个普通关键词加入 trie。
//
// 空字符串会被忽略。
func (n *DFANode) AddWord(word string) {
	if word == "" {
		return
	}
	node := n
	chars := []rune(word)
	for index := range chars {
		node = node.AddChild(chars[index])
	}
	node.End = true
}

// AddChild 返回字符 c 对应的子节点；不存在时会创建。
func (n *DFANode) AddChild(c rune) *DFANode {
	if n.Next == nil {
		n.Next = make(map[rune]*DFANode)
	}

	if next, ok := n.Next[c]; ok {
		return next
	}
	n.Next[c] = &DFANode{
		End:  false,
		Next: nil,
	}
	return n.Next[c]
}

// FindChild 查找字符 c 对应的子节点；不存在时返回 nil。
func (n *DFANode) FindChild(c rune) *DFANode {
	if n.Next == nil {
		return nil
	}

	if next, ok := n.Next[c]; ok {
		return next
	}
	return nil
}
