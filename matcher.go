package matchs

type Matcher interface {
	//Build build Matcher
	Build(words []string)

	//Match return match sensitive words
	Match(text string, onlyOne bool, repl rune) ([]string, string)
}
