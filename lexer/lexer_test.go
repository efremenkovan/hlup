package lexer

import (
	"testing"

	"github.com/efremenkovan/hlup/lang"
	"github.com/efremenkovan/hlup/options"
	"github.com/efremenkovan/hlup/span"
	"github.com/stretchr/testify/require"
)

func Test_Consume(t *testing.T) {
	t.Run("should return token", func(t *testing.T) {
		input := "one"
		lexer := NewLexer(input)

		token, ok := lexer.Consume()
		require.True(t, ok)
		require.Equal(t, token, NewLiteralToken("one", span.NewSpan(0, 2)))
	})

	t.Run("should advance cursor position to next non whitespace char", func(t *testing.T) {
		input := "one two"
		lexer := NewLexer(input)

		token, ok := lexer.Consume()
		require.True(t, ok)
		require.Equal(t, token, NewLiteralToken("one", span.NewSpan(0, 2)))
		require.Equal(t, 4, lexer.cur.position)
	})

	t.Run("should mutate context", func(t *testing.T) {
		input := "'two'"
		lexer := NewLexer(input)

		token, ok := lexer.Consume()
		require.True(t, ok)
		require.Equal(t, token, NewToken(TokenKindQuote, span.NewSpan(0, 0)))
		require.Equal(t, ctxInQuote, lexer.ctx.state)
	})

	t.Run("should skip trailing whitespaces", func(t *testing.T) {
		input := "one      "
		lexer := NewLexer(input)

		token, ok := lexer.Consume()
		require.True(t, ok)
		require.Equal(t, token, NewLiteralToken("one", span.NewSpan(0, 2)))
		require.Equal(t, 9, lexer.cur.position)
		require.True(t, lexer.IsFinished())
	})
}

func Test_Peek(t *testing.T) {
	t.Run("should return token", func(t *testing.T) {
		input := "one"
		lexer := NewLexer(input)

		token, ok := lexer.Peek()
		require.True(t, ok)
		require.Equal(t, token, NewLiteralToken("one", span.NewSpan(0, 2)))
	})

	t.Run("should not advance cursor position to next non whitespace char", func(t *testing.T) {
		input := "one two"
		lexer := NewLexer(input)

		token, ok := lexer.Peek()
		require.True(t, ok)
		require.Equal(t, token, NewLiteralToken("one", span.NewSpan(0, 2)))
		require.Equal(t, 0, lexer.cur.position)
	})

	t.Run("should not mutate context", func(t *testing.T) {
		input := "'two'"
		lexer := NewLexer(input)

		token, ok := lexer.Peek()
		require.True(t, ok)
		require.Equal(t, token, NewToken(TokenKindQuote, span.NewSpan(0, 0)))
		require.Equal(t, int8(0), lexer.ctx.state)
	})
}

func Test_Lexer(t *testing.T) {
	t.Run("literals", func(t *testing.T) {
		t.Run("should consume literal", func(t *testing.T) {
			input := "one"
			lexer := NewLexer(input)

			tokenStream := lexer.TokenStream()
			require.Equal(t, TokenStream{NewLiteralToken("one", span.NewSpan(0, 2))}, tokenStream)
		})

		t.Run("should consume multiword literal", func(t *testing.T) {
			input := "'one two'"
			lexer := NewLexer(input)

			tokenStream := lexer.TokenStream()
			require.Equal(t, TokenStream{NewToken(TokenKindQuote, span.NewSpan(0, 0)), NewLiteralToken("one", span.NewSpan(1, 3)), NewLiteralToken("two", span.NewSpan(5, 7)), NewToken(TokenKindQuote, span.NewSpan(8, 8))}, tokenStream)
		})
	})

	t.Run("keywords", func(t *testing.T) {
		t.Run("en", func(t *testing.T) {
			type test struct {
				name     string
				input    string
				expected TokenStream
			}

			tests := []test{
				{name: "and", input: "one and two", expected: TokenStream{NewLiteralToken("one", span.NewSpan(0, 2)), NewToken(TokenKindKeywordAND, span.NewSpan(4, 6)), NewLiteralToken("two", span.NewSpan(8, 10))}},
				{name: "or", input: "one or two", expected: TokenStream{NewLiteralToken("one", span.NewSpan(0, 2)), NewToken(TokenKindKeywordOR, span.NewSpan(4, 5)), NewLiteralToken("two", span.NewSpan(7, 9))}},
				{name: "not", input: "not one", expected: TokenStream{NewToken(TokenKindKeywordNOT, span.NewSpan(0, 2)), NewLiteralToken("one", span.NewSpan(4, 6))}},
				{name: "all at once", input: "one and two or not three", expected: TokenStream{NewLiteralToken("one", span.NewSpan(0, 2)), NewToken(TokenKindKeywordAND, span.NewSpan(4, 6)), NewLiteralToken("two", span.NewSpan(8, 10)), NewToken(TokenKindKeywordOR, span.NewSpan(12, 13)), NewToken(TokenKindKeywordNOT, span.NewSpan(15, 17)), NewLiteralToken("three", span.NewSpan(19, 23))}},
			}

			for _, tt := range tests {
				t.Run("should parse", func(t *testing.T) {
					lexer := NewLexer(tt.input)
					tokenStream := lexer.TokenStream()
					require.Equal(t, tt.expected, tokenStream)
				})
			}
		})

		t.Run("ru", func(t *testing.T) {
			type test struct {
				name     string
				input    string
				expected TokenStream
			}

			tests := []test{
				{name: "and", input: "one и two", expected: TokenStream{NewLiteralToken("one", span.NewSpan(0, 2)), NewToken(TokenKindKeywordAND, span.NewSpan(4, 4)), NewLiteralToken("two", span.NewSpan(6, 8))}},
				{name: "or", input: "one или two", expected: TokenStream{NewLiteralToken("one", span.NewSpan(0, 2)), NewToken(TokenKindKeywordOR, span.NewSpan(4, 6)), NewLiteralToken("two", span.NewSpan(8, 10))}},
				{name: "not", input: "не one", expected: TokenStream{NewToken(TokenKindKeywordNOT, span.NewSpan(0, 1)), NewLiteralToken("one", span.NewSpan(3, 5))}},
				{name: "all at once", input: "one и two или не three", expected: TokenStream{NewLiteralToken("one", span.NewSpan(0, 2)), NewToken(TokenKindKeywordAND, span.NewSpan(4, 4)), NewLiteralToken("two", span.NewSpan(6, 8)), NewToken(TokenKindKeywordOR, span.NewSpan(10, 12)), NewToken(TokenKindKeywordNOT, span.NewSpan(14, 15)), NewLiteralToken("three", span.NewSpan(17, 21))}},
			}

			for _, tt := range tests {
				t.Run("should parse", func(t *testing.T) {
					lexer := NewLexer(tt.input, options.WithLang(lang.LangRU))
					tokenStream := lexer.TokenStream()
					require.Equal(t, tt.expected, tokenStream)
				})
			}
		})
	})

	t.Run("parentheses", func(t *testing.T) {
		type test struct {
			name     string
			input    string
			expected TokenStream
		}

		tests := []test{
			{name: "empty", input: "()", expected: TokenStream{NewToken(TokenKindLPar, span.NewSpan(0, 0)), NewToken(TokenKindRPar, span.NewSpan(1, 1))}},
			{name: "empty with spaces", input: "( )", expected: TokenStream{NewToken(TokenKindLPar, span.NewSpan(0, 0)), NewToken(TokenKindRPar, span.NewSpan(2, 2))}},
			{name: "nested", input: "(())", expected: TokenStream{NewToken(TokenKindLPar, span.NewSpan(0, 0)), NewToken(TokenKindLPar, span.NewSpan(1, 1)), NewToken(TokenKindRPar, span.NewSpan(2, 2)), NewToken(TokenKindRPar, span.NewSpan(3, 3))}},
			{name: "nested spaced", input: "( ( ) )", expected: TokenStream{NewToken(TokenKindLPar, span.NewSpan(0, 0)), NewToken(TokenKindLPar, span.NewSpan(2, 2)), NewToken(TokenKindRPar, span.NewSpan(4, 4)), NewToken(TokenKindRPar, span.NewSpan(6, 6))}},
			{name: "nested with content", input: "((one))", expected: TokenStream{NewToken(TokenKindLPar, span.NewSpan(0, 0)), NewToken(TokenKindLPar, span.NewSpan(1, 1)), NewLiteralToken("one", span.NewSpan(2, 4)), NewToken(TokenKindRPar, span.NewSpan(5, 5)), NewToken(TokenKindRPar, span.NewSpan(6, 6))}},
			{name: "nested spaced", input: "(one and two) or not (three and four)", expected: TokenStream{
				NewToken(TokenKindLPar, span.NewSpan(0, 0)),
				NewLiteralToken("one", span.NewSpan(1, 3)),
				NewToken(TokenKindKeywordAND, span.NewSpan(5, 7)),
				NewLiteralToken("two", span.NewSpan(9, 11)),
				NewToken(TokenKindRPar, span.NewSpan(12, 12)),
				NewToken(TokenKindKeywordOR, span.NewSpan(14, 15)),
				NewToken(TokenKindKeywordNOT, span.NewSpan(17, 19)),
				NewToken(TokenKindLPar, span.NewSpan(21, 21)),
				NewLiteralToken("three", span.NewSpan(22, 26)),
				NewToken(TokenKindKeywordAND, span.NewSpan(28, 30)),
				NewLiteralToken("four", span.NewSpan(32, 35)),
				NewToken(TokenKindRPar, span.NewSpan(36, 36)),
			}},
		}

		for _, tt := range tests {
			t.Run("should parse", func(t *testing.T) {
				lexer := NewLexer(tt.input)
				tokenStream := lexer.TokenStream()
				require.Equal(t, tt.expected, tokenStream)
			})
		}
	})
}
