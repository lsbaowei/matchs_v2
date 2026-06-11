package matchs

import "sort"

// DFAMatcher 使用 Aho-Corasick 自动机匹配普通关键词。
//
// Build 阶段会构造 trie 和 failure 指针；Match 阶段单次扫描文本并收集命中结果。
// 替换逻辑按 rune 下标执行，适用于中文等多字节字符。输出顺序保持为按命中起点优先，
// 同一起点下按关键词长度递增。
type DFAMatcher struct {
	root          *DFANode
	maxWordLength int
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

// Build 构造 Aho-Corasick 自动机。
//
// 空字符串关键词会被忽略，避免根节点被误标记为命中节点。
func (d *DFAMatcher) Build(words []string) {
	for _, item := range words {
		if length := len([]rune(item)); length > d.maxWordLength {
			d.maxWordLength = length
		}
		d.root.AddWord(item)
	}
	d.buildFailure()
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

	if onlyOne {
		match, ok := d.findFirstMatch(textChars)
		if !ok {
			replaceText = text
			return
		}
		sensitiveWords = append(sensitiveWords, match.word)
		replaceRune(textCharsCopy, repl, match.begin, match.end)
		replaceText = string(textCharsCopy)
		return
	}

	matches := d.findMatches(textChars)
	if len(matches) == 0 {
		replaceText = text
		return
	}

	for _, match := range matches {
		sensitiveWords = append(sensitiveWords, match.word)
		replaceRune(textCharsCopy, repl, match.begin, match.end)
	}

	replaceText = string(textCharsCopy)
	return
}

func (d *DFAMatcher) findFirstMatch(textChars []rune) (dfaMatch, bool) {
	node := d.root
	var best dfaMatch
	found := false

	for index, c := range textChars {
		node = d.nextNode(node, c)
		for outputNode := node; outputNode != nil && outputNode != d.root; outputNode = outputNode.fail {
			for _, output := range outputNode.outputs {
				begin := index - output.length + 1
				if begin < 0 {
					continue
				}
				match := dfaMatch{
					word:  output.word,
					begin: begin,
					end:   index + 1,
				}
				if !found || lessDFAMatch(match, best) {
					best = match
					found = true
				}
			}
		}
		if found && index-best.begin+1 >= d.maxWordLength {
			return best, true
		}
	}

	return best, found
}

func (d *DFAMatcher) buildFailure() {
	if d.root == nil {
		return
	}

	d.root.fail = d.root
	queue := make([]*DFANode, 0)
	for _, child := range d.root.Next {
		child.fail = d.root
		queue = append(queue, child)
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for c, child := range node.Next {
			fail := node.fail
			for fail != d.root && fail.FindChild(c) == nil {
				fail = fail.fail
			}
			if next := fail.FindChild(c); next != nil && next != child {
				child.fail = next
			} else {
				child.fail = d.root
			}
			queue = append(queue, child)
		}
	}
}

func (d *DFAMatcher) findMatches(textChars []rune) []dfaMatch {
	node := d.root
	matches := make([]dfaMatch, 0)

	for index, c := range textChars {
		node = d.nextNode(node, c)

		for outputNode := node; outputNode != nil && outputNode != d.root; outputNode = outputNode.fail {
			for _, output := range outputNode.outputs {
				begin := index - output.length + 1
				if begin < 0 {
					continue
				}
				matches = append(matches, dfaMatch{
					word:  output.word,
					begin: begin,
					end:   index + 1,
				})
			}
		}
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return lessDFAMatch(matches[i], matches[j])
	})

	return matches
}

func (d *DFAMatcher) nextNode(node *DFANode, c rune) *DFANode {
	for node != d.root && node.FindChild(c) == nil {
		node = node.fail
	}
	if next := node.FindChild(c); next != nil {
		return next
	}
	return d.root
}

func lessDFAMatch(left, right dfaMatch) bool {
	if left.begin != right.begin {
		return left.begin < right.begin
	}
	return left.end < right.end
}

// replaceRune 将 [begin, end) 范围内的 rune 替换为 replaceChar。
func replaceRune(chars []rune, replaceChar rune, begin int, end int) {
	for i := begin; i < end; i++ {
		chars[i] = replaceChar
	}
}

// DFANode 是 DFA trie 的节点。
type DFANode struct {
	End     bool
	Next    map[rune]*DFANode
	fail    *DFANode
	outputs []dfaOutput
}

type dfaOutput struct {
	word   string
	length int
}

type dfaMatch struct {
	word  string
	begin int
	end   int
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
	for _, output := range node.outputs {
		if output.word == word {
			return
		}
	}
	node.End = true
	node.outputs = append(node.outputs, dfaOutput{
		word:   word,
		length: len(chars),
	})
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
