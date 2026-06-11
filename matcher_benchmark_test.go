package matchs

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkDFAMatcherBuild(b *testing.B) {
	for _, size := range []int{1000, 10000, 100000} {
		words := benchmarkWords(size)
		b.Run(fmt.Sprintf("words_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				matcher := NewDFAMatcher()
				matcher.Build(words)
			}
		})
	}
}

func BenchmarkDFAMatcherMatch(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		matcher := NewDFAMatcher()
		matcher.Build(benchmarkWords(size))
		text := strings.Repeat("这是一段普通文本", 100) + fmt.Sprintf("关键词%06d", size-1)

		b.Run(fmt.Sprintf("words_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				matcher.Match(text, false, '*')
			}
		})
	}
}

func BenchmarkDFAMatcherMatchOnlyOne(b *testing.B) {
	matcher := NewDFAMatcher()
	matcher.Build(benchmarkWords(10000))
	text := "关键词000001" + strings.Repeat("这是一段普通文本", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher.Match(text, true, '*')
	}
}

func BenchmarkDFAMatcherCommonPrefix(b *testing.B) {
	matcher := NewDFAMatcher()
	matcher.Build(benchmarkCommonPrefixWords(10000))
	text := strings.Repeat("公共前缀", 100) + "公共前缀09999"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher.Match(text, false, '*')
	}
}

func BenchmarkAssembleMatcherMatch(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		matcher := NewAssembleMatcher()
		matcher.Build(benchmarkAssembleRules(size))
		text := fmt.Sprintf("A%06d xx B%06d", size-1, size-1)

		b.Run(fmt.Sprintf("rules_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				matcher.Match(text, false, '*')
			}
		})
	}
}

func BenchmarkRegexpMatcherMatch(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		matcher := NewRegexpMatcher()
		matcher.Build(benchmarkRegexpRules(size))
		text := fmt.Sprintf("prefix%06d123", size-1)

		b.Run(fmt.Sprintf("rules_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				matcher.Match(text, false, '*')
			}
		})
	}
}

func benchmarkWords(size int) []string {
	words := make([]string, 0, size)
	for i := 0; i < size; i++ {
		words = append(words, fmt.Sprintf("关键词%06d", i))
	}
	return words
}

func benchmarkCommonPrefixWords(size int) []string {
	words := make([]string, 0, size)
	for i := 0; i < size; i++ {
		words = append(words, fmt.Sprintf("公共前缀%05d", i))
	}
	return words
}

func benchmarkAssembleRules(size int) []string {
	rules := make([]string, 0, size)
	for i := 0; i < size; i++ {
		rules = append(rules, fmt.Sprintf("A%06d|B%06d", i, i))
	}
	return rules
}

func benchmarkRegexpRules(size int) []string {
	rules := make([]string, 0, size)
	for i := 0; i < size; i++ {
		rules = append(rules, fmt.Sprintf(`prefix%06d\d+`, i))
	}
	return rules
}
