package expression

import (
	"testing"
	"unicode/utf8"

	"github.com/efremenkovan/hlup/span"
	"github.com/stretchr/testify/require"
)

// Generate a token sequence as if each token is separated by space
func ts(tokens ...string) TokenStream {
	t := make(TokenStream, 0)
	spanCur := 0
	for _, token := range tokens {
		t = append(t, Token{Value: token, Span: span.NewSpan(spanCur, spanCur+utf8.RuneCountInString(token)-1)})
		spanCur += utf8.RuneCountInString(token) + 1
	}

	return t
}

func Test_TokenSequence_ExactMatch(t *testing.T) {
	t.Run("should match", func(t *testing.T) {
		input := ts("previous", "some", "long", "inputs", "following")
		expr := ts("some", "long", "inputs")

		matched, s := expr.MatchEqual(input)
		require.True(t, matched)
		require.Len(t, s, 3)
		require.Equal(t, []span.Span{{Start: 9, End: 12}, {Start: 14, End: 17}, {Start: 19, End: 24}}, s)
	})

	t.Run("should not match", func(t *testing.T) {
		type test struct {
			name  string
			input TokenStream
			expr  TokenStream
		}

		tests := []test{
			{
				name:  "missing last element",
				input: ts("previous", "some", "long", "following"),
				expr:  ts("some", "long", "inputs"),
			},
			{
				name:  "missing central element",
				input: ts("previous", "some", "inputs", "following"),
				expr:  ts("some", "long", "inputs"),
			},
			{
				name:  "extra token in between matching ones",
				input: ts("previous", "some", "long", "long", "inputs", "following"),
				expr:  ts("some", "long", "inputs"),
			},
			{
				name:  "some token contains but not equal",
				input: ts("previous", "some", "oolongoo", "inputs", "following"),
				expr:  ts("some", "long", "inputs"),
			},
		}

		for _, tt := range tests {
			matched, s := tt.expr.MatchEqual(tt.input)
			require.False(t, matched)
			require.Len(t, s, 0)
		}
	})
}

func Test_TokenSequence_RegularMatch(t *testing.T) {
	t.Run("should match", func(t *testing.T) {
		input := ts("previous", "some", "long", "inputs", "following")
		expr := ts("some", "long", "inputs")

		matched, s := expr.MatchContains(input)
		require.True(t, matched)
		require.Len(t, s, 3)
		require.Equal(t, []span.Span{
			{Start: 9, End: 12},
			{Start: 14, End: 17},
			{Start: 19, End: 24},
		}, s)
	})

	t.Run("should match containing", func(t *testing.T) {
		input := ts("previous", "some", "oolongoo", "inputs", "following")
		expr := ts("some", "long", "inputs")

		matched, s := expr.MatchContains(input)
		require.True(t, matched)
		require.Len(t, s, 3)
		require.Equal(t, []span.Span{
			{Start: 9, End: 12},
			{Start: 14, End: 21},
			{Start: 23, End: 28},
		}, s)
	})

	t.Run("should not match", func(t *testing.T) {
		type test struct {
			name  string
			input TokenStream
			expr  TokenStream
		}

		tests := []test{
			{
				name:  "missing last element",
				input: ts("previous", "some", "long", "following"),
				expr:  ts("some", "long", "inputs"),
			},
			{
				name:  "missing central element",
				input: ts("previous", "some", "inputs", "following"),
				expr:  ts("some", "long", "inputs"),
			},
			{
				name:  "extra token in between matching ones",
				input: ts("previous", "some", "long", "long", "inputs", "following"),
				expr:  ts("some", "long", "inputs"),
			},
		}

		for _, tt := range tests {
			matched, s := tt.expr.MatchContains(tt.input)
			require.False(t, matched)
			require.Len(t, s, 0)
		}
	})
}
