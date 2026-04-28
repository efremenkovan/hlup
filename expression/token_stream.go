package expression

import (
	"strings"

	"github.com/efremenkovan/hlup/span"
)

type Token struct {
	Value string
	Span  span.Span
}

func (t Token) Contains(rhs Token) bool {
	return strings.Contains(string(t.Value), string(rhs.Value))
}

// TokenStream is a sequence of words that can be matched against the expression
type TokenStream []Token

func matchEqualFunc(actual, expected Token) bool {
	return actual.Value == expected.Value
}

func matchContainsFunc(actual, expected Token) bool {
	return strings.Contains(actual.Value, expected.Value)
}

func (s TokenStream) MatchEqual(input TokenStream) (bool, []span.Span) {
	return s.match(input, matchEqualFunc)
}

func (s TokenStream) MatchContains(input TokenStream) (bool, []span.Span) {
	return s.match(input, matchContainsFunc)
}

func (s TokenStream) match(input TokenStream, matchFunc func(actual, expected Token) bool) (bool, []span.Span) {
	// Empty expected token sequence is equal to any other token sequence
	if len(s) == 0 {
		return true, nil
	}

	if len(input) == 0 {
		return false, nil
	}

	matchedTokensCount := 0
	matchedSpans := make([]span.Span, len(s))

	// We will come back to this point in case potential sequence match will fail before reaching the end
	persistedSequenceMatchStartIndex := -1
	inputIndex := 0
	for inputIndex < len(input) {
		inputToken := input[inputIndex]

		if !matchFunc(inputToken, s[matchedTokensCount]) {
			matchedTokensCount = 0

			if persistedSequenceMatchStartIndex != -1 {
				// Go to the rune following the previous failed sequence match start
				inputIndex = persistedSequenceMatchStartIndex + 1
				persistedSequenceMatchStartIndex = -1
				continue
			}

			inputIndex++
			continue
		}

		if matchedTokensCount == 0 {
			persistedSequenceMatchStartIndex = inputIndex
		}

		matchedSpans[matchedTokensCount] = inputToken.Span
		matchedTokensCount += 1

		if matchedTokensCount == len(s) {
			return true, matchedSpans
		}

		inputIndex++
	}

	return false, nil
}
