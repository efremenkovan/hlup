package hlup

import (
	"testing"

	"github.com/efremenkovan/hlup/expression"
	"github.com/efremenkovan/hlup/span"
	"github.com/stretchr/testify/require"
)

func Test_TokenizeOptions(t *testing.T) {
	getDefaultOpts := func() tokenizeOptions {
		return tokenizeOptions{
			replaceTable:      map[rune]rune{'a': 'b'},
			wordBreakersTable: map[rune]struct{}{'a': {}},
		}
	}

	type testCase struct {
		name     string
		funcs    []tokenizeOptionsPatchFunc
		expected tokenizeOptions
	}

	tests := []testCase{
		{
			name:  "extend replace table",
			funcs: []tokenizeOptionsPatchFunc{WithExtendedReplaceTable(map[rune]rune{'b': 'c'})},
			expected: tokenizeOptions{
				replaceTable:      map[rune]rune{'a': 'b', 'b': 'c'},
				wordBreakersTable: map[rune]struct{}{'a': {}},
			},
		},
		{
			name:  "custom replace table",
			funcs: []tokenizeOptionsPatchFunc{WithCustomReplaceTable(map[rune]rune{'b': 'c'})},
			expected: tokenizeOptions{
				replaceTable:      map[rune]rune{'b': 'c'},
				wordBreakersTable: map[rune]struct{}{'a': {}},
			},
		},
		{
			name:  "extend word breakers table",
			funcs: []tokenizeOptionsPatchFunc{WithExtendedWordBreakersList([]rune{'b'})},
			expected: tokenizeOptions{
				replaceTable:      map[rune]rune{'a': 'b'},
				wordBreakersTable: map[rune]struct{}{'a': {}, 'b': {}},
			},
		},
		{
			name:  "custom work breakers table",
			funcs: []tokenizeOptionsPatchFunc{WithCustomWordBreakersList([]rune{'b'})},
			expected: tokenizeOptions{
				replaceTable:      map[rune]rune{'a': 'b'},
				wordBreakersTable: map[rune]struct{}{'b': {}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initial := getDefaultOpts()
			for _, f := range tt.funcs {
				f(&initial)
			}

			require.Equal(t, tt.expected, initial)
		})
	}
}

func Test_TokenizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		optFuncs []tokenizeOptionsPatchFunc
		want     expression.TokenStream
	}{
		{
			name:  "single word",
			input: "single",
			want: expression.TokenStream{
				expression.Token{Value: "single", Span: span.NewSpan(0, 5)},
			},
		},
		{
			name:  "replace ё with e",
			input: "ёлка",
			want: expression.TokenStream{
				expression.Token{Value: "елка", Span: span.NewSpan(0, 3)},
			},
		},
		{
			name:  "remove punctuation",
			input: "these are. some really, really 'separate' sentences!",
			want: expression.TokenStream{
				expression.Token{Value: "these", Span: span.NewSpan(0, 4)},
				expression.Token{Value: "are", Span: span.NewSpan(6, 8)},
				expression.Token{Value: "some", Span: span.NewSpan(11, 14)},
				expression.Token{Value: "really", Span: span.NewSpan(16, 21)},
				expression.Token{Value: "really", Span: span.NewSpan(24, 29)},
				expression.Token{Value: "separate", Span: span.NewSpan(32, 39)},
				expression.Token{Value: "sentences", Span: span.NewSpan(42, 50)},
			},
		},
		{
			name:  "lowercase",
			input: "Some OF tHe WordS",

			want: expression.TokenStream{
				expression.Token{Value: "some", Span: span.NewSpan(0, 3)},
				expression.Token{Value: "of", Span: span.NewSpan(5, 6)},
				expression.Token{Value: "the", Span: span.NewSpan(8, 10)},
				expression.Token{Value: "words", Span: span.NewSpan(12, 16)},
			},
		},
		{
			name:     "with custom replace set",
			input:    "ёлка",
			optFuncs: []tokenizeOptionsPatchFunc{WithCustomReplaceTable(map[rune]rune{'л': 'м', 'к': 'а', 'а': 'е'})},
			want: expression.TokenStream{
				expression.Token{Value: "ёмае", Span: span.NewSpan(0, 3)},
			},
		},
		{
			name:     "with custom break chars",
			input:    "found the fire",
			optFuncs: []tokenizeOptionsPatchFunc{WithCustomWordBreakersList([]rune{'f'})},
			want: expression.TokenStream{
				expression.Token{Value: "ound the ", Span: span.NewSpan(1, 9)},
				expression.Token{Value: "ire", Span: span.NewSpan(11, 13)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TokenizeInput(tt.input, tt.optFuncs...)
			if !tokenSequenceEqual(got, tt.want) {
				t.Errorf("tokenize(%q):\nvalues: %v\nwant:   %v\n\nspans: %v\nwant:  %v", tt.input, got, tt.want, spans(got), spans(tt.want))
			}
		})
	}
}

func Test_Hlup(t *testing.T) {
	t.Run("should match equal an expression in text", func(t *testing.T) {
		type testCase struct {
			name        string
			input       string
			shouldMatch bool
			span        []span.Span
		}

		rule := "go and (lang or way or idiomatic) and not 'to hell'"
		tests := []testCase{
			{
				name:        "direct sequence",
				input:       "this is a go way",
				shouldMatch: true,
				span:        []span.Span{span.NewSpan(10, 11), span.NewSpan(13, 15)},
			},
			{
				name:        "indirect sequence",
				input:       "idiomatic go",
				shouldMatch: true,
				span:        []span.Span{span.NewSpan(10, 11), span.NewSpan(0, 8)},
			},
			{
				name:        "with part of restricted word",
				input:       "idiomatic go hell",
				shouldMatch: true,
				span:        []span.Span{span.NewSpan(10, 11), span.NewSpan(0, 8)},
			},
			{
				name:        "with restricted word",
				input:       "idiomatic go to hell",
				shouldMatch: false,
			},
			{
				name:        "with no matching sequence",
				input:       "idiomatic hell",
				shouldMatch: false,
			},
			{
				name:        "contains but not equal",
				input:       "this is nogo way",
				shouldMatch: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				e, err := CompileExpression(rule)
				require.NoError(t, err)
				isMatch, foundSpans := e.MatchEqual(TokenizeInput(tt.input))
				require.Equal(t, tt.shouldMatch, isMatch)
				if tt.shouldMatch {
					require.Equal(t, tt.span, foundSpans)
				}
			})
		}
	})

	t.Run("should match equal an expression in text", func(t *testing.T) {
		type testCase struct {
			name        string
			input       string
			shouldMatch bool
			span        []span.Span
		}

		rule := "go and (lang or way or idiomatic) and not 'to hell'"
		tests := []testCase{
			{
				name:        "direct sequence",
				input:       "this is a go way",
				shouldMatch: true,
				span:        []span.Span{span.NewSpan(10, 11), span.NewSpan(13, 15)},
			},
			{
				name:        "indirect sequence",
				input:       "idiomatic go",
				shouldMatch: true,
				span:        []span.Span{span.NewSpan(10, 11), span.NewSpan(0, 8)},
			},
			{
				name:        "with part of restricted word",
				input:       "idiomatic go  hell",
				shouldMatch: true,
				span:        []span.Span{span.NewSpan(10, 11), span.NewSpan(0, 8)},
			},
			{
				name:        "with restricted word",
				input:       "idiomatic go to hell",
				shouldMatch: false,
			},
			{
				name:        "with no matching sequence",
				input:       "idiomatic hell",
				shouldMatch: false,
			},
			{
				name:        "contains but not equal",
				input:       "this is nogo way",
				shouldMatch: true,
				span:        []span.Span{span.NewSpan(8, 11), span.NewSpan(13, 15)},
			},
			{
				name:        "every contains",
				input:       "oooidiomaticooo vvvvgovvv",
				shouldMatch: true,
				span:        []span.Span{span.NewSpan(16, 24), span.NewSpan(0, 14)},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				e, err := CompileExpression(rule)
				require.NoError(t, err)
				isMatch, foundSpans := e.MatchContains(TokenizeInput(tt.input))
				require.Equal(t, tt.shouldMatch, isMatch)
				if tt.shouldMatch {
					require.Equal(t, tt.span, foundSpans)
				}
			})
		}
	})
}

func spans(seq expression.TokenStream) []span.Span {
	res := make([]span.Span, len(seq))
	for i, t := range seq {
		res[i] = t.Span
	}

	return res
}

func tokenSequenceEqual(a, b expression.TokenStream) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Value != b[i].Value {
			return false
		}
		if a[i].Span != b[i].Span {
			return false
		}
	}
	return true
}
