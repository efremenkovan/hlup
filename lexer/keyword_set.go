package lexer

import (
	"github.com/efremenkovan/hlup/lang"
)

// Set of string representations of all the tokens supported by hlup for a specific language
type keywordSet struct {
	lang lang.Lang
	set  map[string]TokenKind
}

var ruKeywordSet = keywordSet{
	lang: lang.LangRU,
	set: map[string]TokenKind{
		"и":   TokenKindKeywordAND,
		"или": TokenKindKeywordOR,
		"не":  TokenKindKeywordNOT,
	},
}

var engKeywordSet = keywordSet{
	lang: lang.LangEN,
	set: map[string]TokenKind{
		"and": TokenKindKeywordAND,
		"or":  TokenKindKeywordOR,
		"not": TokenKindKeywordNOT,
	},
}

// KeywordTokenKind returns TokenKind (if any) represented by input string in this set
func (s keywordSet) KeywordTokenKind(input string) (TokenKind, bool) {
	v, ok := s.set[input]
	return v, ok
}

// Returns keyword for the provided language
func getSetByLang(lng lang.Lang) *keywordSet {
	switch lng {
	case lang.LangRU:
		return &ruKeywordSet
	case lang.LangEN:
		return &engKeywordSet
	}

	// Fallback to english if provided with unknown language
	return &engKeywordSet
}
