package matchs_v2

import (
	"reflect"
	"testing"
)

func TestDFAMatcherMatchAndReplace(t *testing.T) {
	matcher := NewDFAMatcher()
	matcher.Build([]string{"坏", "坏人", ""})

	words, replaced := matcher.Match("坏人来了", false, '*')

	wantWords := []string{"坏", "坏人"}
	if !reflect.DeepEqual(words, wantWords) {
		t.Fatalf("words = %#v, want %#v", words, wantWords)
	}
	if replaced != "**来了" {
		t.Fatalf("replaced = %q, want %q", replaced, "**来了")
	}
}

func TestDFAMatcherOnlyOne(t *testing.T) {
	matcher := NewDFAMatcher()
	matcher.Build([]string{"敏感词"})

	words, replaced := matcher.Match("这里有敏感词和敏感词", true, '*')

	wantWords := []string{"敏感词"}
	if !reflect.DeepEqual(words, wantWords) {
		t.Fatalf("words = %#v, want %#v", words, wantWords)
	}
	if replaced != "这里有***和敏感词" {
		t.Fatalf("replaced = %q, want %q", replaced, "这里有***和敏感词")
	}
}

func TestDFAMatcherOverlappingKeepsStartOrder(t *testing.T) {
	matcher := NewDFAMatcher()
	matcher.Build([]string{"abc", "b", "bc"})

	words, replaced := matcher.Match("abc", false, '*')

	wantWords := []string{"abc", "b", "bc"}
	if !reflect.DeepEqual(words, wantWords) {
		t.Fatalf("words = %#v, want %#v", words, wantWords)
	}
	if replaced != "***" {
		t.Fatalf("replaced = %q, want %q", replaced, "***")
	}
}

func TestAssembleMatcherANDOrder(t *testing.T) {
	matcher := NewAssembleMatcher()
	matcher.Build([]string{"A|B|C"})

	words, _ := matcher.Match("AxxBxxC", false, '*')
	if !reflect.DeepEqual(words, []string{"A|B|C"}) {
		t.Fatalf("ordered words = %#v, want %#v", words, []string{"A|B|C"})
	}

	words, _ = matcher.Match("BxxAxxC", false, '*')
	if len(words) != 0 {
		t.Fatalf("unordered words = %#v, want empty", words)
	}
}

func TestAssembleMatcherPrefilter(t *testing.T) {
	matcher := NewAssembleMatcher()
	matcher.Build([]string{"A|B|C", "X|Y"})

	if _, ok := matcher.anchorIndex["A"]; !ok {
		t.Fatal("anchor A not found")
	}
	if _, ok := matcher.anchorIndex["X"]; !ok {
		t.Fatal("anchor X not found")
	}

	words, _ := matcher.Match("xxBxxC", false, '*')
	if len(words) != 0 {
		t.Fatalf("words = %#v, want empty", words)
	}
}

func TestAssembleMatcherNON(t *testing.T) {
	matcher := NewAssembleMatcher()
	matcher.Build([]string{"A#B"})

	words, replaced := matcher.Match("xxAxx", false, '*')
	if !reflect.DeepEqual(words, []string{"A#B"}) {
		t.Fatalf("words = %#v, want %#v", words, []string{"A#B"})
	}
	if replaced != "xxAxx" {
		t.Fatalf("replaced = %q, want %q", replaced, "xxAxx")
	}

	words, _ = matcher.Match("xxAxxBxx", false, '*')
	if len(words) != 0 {
		t.Fatalf("excluded words = %#v, want empty", words)
	}
}

func TestRegexpMatcher(t *testing.T) {
	matcher := NewRegexpMatcher()
	matcher.Build([]string{`1\d{10}`})

	words, replaced := matcher.Match("手机号13800138000", false, '*')

	if !reflect.DeepEqual(words, []string{`reg@1\d{10}`}) {
		t.Fatalf("words = %#v, want %#v", words, []string{`reg@1\d{10}`})
	}
	if replaced != "手机号13800138000" {
		t.Fatalf("replaced = %q, want %q", replaced, "手机号13800138000")
	}
}

func TestRegexpMatcherPrefilterAndFallback(t *testing.T) {
	matcher := NewRegexpMatcher()
	matcher.Build([]string{`hello\d+`, `\d{11}`})

	if _, ok := matcher.anchorIndex["hello"]; !ok {
		t.Fatal("anchor hello not found")
	}
	if len(matcher.fallbackRules) != 1 {
		t.Fatalf("fallbackRules length = %d, want 1", len(matcher.fallbackRules))
	}

	words, _ := matcher.Match("手机号13800138000", false, '*')
	if !reflect.DeepEqual(words, []string{`reg@\d{11}`}) {
		t.Fatalf("words = %#v, want %#v", words, []string{`reg@\d{11}`})
	}
}

func TestMatchServiceReturnsReplacement(t *testing.T) {
	service := NewMatchService()
	service.Build([]string{"敏感词", "A|B", `reg@1\d{10}`})

	words, replaced := service.Match("敏感词 AxxB 手机号13800138000", false, '*')

	wantWords := []string{"敏感词", "A|B", `reg@1\d{10}`}
	if !reflect.DeepEqual(words, wantWords) {
		t.Fatalf("words = %#v, want %#v", words, wantWords)
	}
	if replaced != "*** AxxB 手机号13800138000" {
		t.Fatalf("replaced = %q, want %q", replaced, "*** AxxB 手机号13800138000")
	}
}
