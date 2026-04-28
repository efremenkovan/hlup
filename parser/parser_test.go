package parser

import (
	"fmt"
	"testing"

	"github.com/efremenkovan/hlup/expression"
	"github.com/efremenkovan/hlup/lexer"
	"github.com/efremenkovan/hlup/span"
	"github.com/stretchr/testify/require"
)

func l(value string, span span.Span) lexer.Token {
	return lexer.NewLiteralToken(value, span)
}

func not(span span.Span) lexer.Token {
	return lexer.NewToken(lexer.TokenKindKeywordNOT, span)
}

func or(span span.Span) lexer.Token {
	return lexer.NewToken(lexer.TokenKindKeywordOR, span)
}

func and(span span.Span) lexer.Token {
	return lexer.NewToken(lexer.TokenKindKeywordAND, span)
}

func lp(span span.Span) lexer.Token {
	return lexer.NewToken(lexer.TokenKindLPar, span)
}

func rp(span span.Span) lexer.Token {
	return lexer.NewToken(lexer.TokenKindRPar, span)
}

func q(span span.Span) lexer.Token {
	return lexer.NewToken(lexer.TokenKindQuote, span)
}

func Test_Parse(t *testing.T) {
	type test struct {
		name        string
		input       lexer.TokenStream
		expected    expression.Expression
		expectedErr error
	}

	tests := []test{
		{
			name: "and",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				l("two", span.NewSpan(8, 10)),
			},
			expected: expression.AndExpression{
				Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(0, 2)}},
				Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(8, 10)}},
			},
		},
		{
			name: "or",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				or(span.NewSpan(4, 5)),
				l("two", span.NewSpan(7, 9)),
			},
			expected: expression.OrExpression{
				Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(0, 2)}},
				Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(7, 9)}},
			},
		},
		{
			name: "root level standalone not",
			input: lexer.TokenStream{
				not(span.NewSpan(0, 2)),
				l("one", span.NewSpan(4, 6)),
			},
			expected: expression.NotExpression{
				Expr: expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(4, 6)}},
			},
		},
		{
			name: "not with parentheses",
			input: lexer.TokenStream{
				not(span.NewSpan(0, 2)),
				lp(span.NewSpan(4, 4)),
				l("one", span.NewSpan(5, 7)),
				and(span.NewSpan(9, 11)),
				l("two", span.NewSpan(13, 15)),
				rp(span.NewSpan(16, 16)),
			},
			expected: expression.NotExpression{
				Expr: expression.AndExpression{
					Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(5, 7)}},
					Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(13, 15)}},
				},
			},
		},
		{
			name: "or with negated left branch",
			input: lexer.TokenStream{
				not(span.NewSpan(0, 2)),
				lp(span.NewSpan(3, 3)),
				l("one", span.NewSpan(4, 6)),
				and(span.NewSpan(8, 10)),
				l("two", span.NewSpan(12, 14)),
				rp(span.NewSpan(15, 15)),
				or(span.NewSpan(17, 18)),
				lp(span.NewSpan(20, 20)),
				l("three", span.NewSpan(21, 25)),
				and(span.NewSpan(27, 29)),
				l("four", span.NewSpan(31, 33)),
				rp(span.NewSpan(34, 34)),
			},
			expected: expression.OrExpression{
				Left: expression.NotExpression{
					Expr: expression.AndExpression{
						Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(4, 6)}},
						Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(12, 14)}},
					},
				},
				Right: expression.AndExpression{
					Left:  expression.TokenStream{expression.Token{Value: "three", Span: span.NewSpan(21, 25)}},
					Right: expression.TokenStream{expression.Token{Value: "four", Span: span.NewSpan(31, 33)}},
				},
			},
		},
		{
			name: "multiple nested levels",
			input: lexer.TokenStream{
				not(span.NewSpan(0, 2)),
				lp(span.NewSpan(3, 3)),
				lp(span.NewSpan(4, 4)),
				l("one", span.NewSpan(5, 7)),
				or(span.NewSpan(9, 10)),
				l("two", span.NewSpan(12, 14)),
				rp(span.NewSpan(15, 15)),
				and(span.NewSpan(17, 19)),
				lp(span.NewSpan(21, 21)),
				l("three", span.NewSpan(22, 26)),
				or(span.NewSpan(28, 29)),
				l("four", span.NewSpan(31, 34)),
				rp(span.NewSpan(36, 36)),
				rp(span.NewSpan(37, 37)),
				or(span.NewSpan(39, 40)),
				lp(span.NewSpan(42, 42)),
				l("five", span.NewSpan(43, 46)),
				and(span.NewSpan(48, 50)),
				l("six", span.NewSpan(52, 54)),
				rp(span.NewSpan(55, 55)),
			},
			expected: expression.OrExpression{
				Left: expression.NotExpression{
					Expr: expression.AndExpression{
						Left: expression.OrExpression{
							Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(5, 7)}},
							Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(12, 14)}},
						},
						Right: expression.OrExpression{
							Left:  expression.TokenStream{expression.Token{Value: "three", Span: span.NewSpan(22, 26)}},
							Right: expression.TokenStream{expression.Token{Value: "four", Span: span.NewSpan(31, 34)}},
						},
					},
				},
				Right: expression.AndExpression{
					Left:  expression.TokenStream{expression.Token{Value: "five", Span: span.NewSpan(43, 46)}},
					Right: expression.TokenStream{expression.Token{Value: "six", Span: span.NewSpan(52, 54)}},
				},
			},
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				lp(span.NewSpan(0, 0)),
				l("one", span.NewSpan(1, 3)),
				and(span.NewSpan(5, 7)),
				l("two", span.NewSpan(9, 11)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxUnbalancedParentheses, span.NewSpan(11, 11)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				lp(span.NewSpan(0, 0)),
				l("one", span.NewSpan(1, 3)),
				and(span.NewSpan(5, 7)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxAndNoRightExpr, span.NewSpan(5, 7)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				not(span.NewSpan(8, 10)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxNotNoFollowingExpr, span.NewSpan(8, 10)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				or(span.NewSpan(4, 5)),
				not(span.NewSpan(7, 9)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxNotNoFollowingExpr, span.NewSpan(7, 9)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				or(span.NewSpan(4, 5)),
				not(span.NewSpan(7, 9)),
				or(span.NewSpan(11, 13)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxNotInvalidFollowing, span.NewSpan(11, 13)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				or(span.NewSpan(4, 5)),
				or(span.NewSpan(7, 8)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxOrInvalidFollowing, span.NewSpan(7, 8)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				or(span.NewSpan(4, 5)),
				and(span.NewSpan(7, 9)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxOrInvalidFollowing, span.NewSpan(7, 9)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				and(span.NewSpan(8, 10)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxAndInvalidFollowing, span.NewSpan(8, 10)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				or(span.NewSpan(8, 9)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxAndInvalidFollowing, span.NewSpan(8, 9)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				l("one", span.NewSpan(8, 10)),
				l("one", span.NewSpan(12, 14)),
			},
			expectedErr: newParserError(ErrInvalidSyntax, span.NewSpan(12, 14)),
		},
		{
			name: "invalid syntax",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				lp(span.NewSpan(4, 4)),
				and(span.NewSpan(5, 7)),
				l("two", span.NewSpan(9, 11)),
				rp(span.NewSpan(12, 12)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxAndNoLeftExpr, span.NewSpan(5, 7)),
		},
		{
			name: "unbalanced parantheses",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				l("two", span.NewSpan(8, 10)),
				rp(span.NewSpan(11, 11)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxUnbalancedParentheses, span.NewSpan(11, 11)),
		},
		{
			name: "chained and is left-associative",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				l("two", span.NewSpan(8, 10)),
				and(span.NewSpan(12, 14)),
				l("three", span.NewSpan(16, 20)),
			},
			expected: expression.AndExpression{
				Left: expression.AndExpression{
					Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(0, 2)}},
					Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(8, 10)}},
				},
				Right: expression.TokenStream{expression.Token{Value: "three", Span: span.NewSpan(16, 20)}},
			},
		},
		{
			name: "chained or is left-associative",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				or(span.NewSpan(4, 5)),
				l("two", span.NewSpan(7, 9)),
				or(span.NewSpan(11, 12)),
				l("three", span.NewSpan(14, 18)),
			},
			expected: expression.OrExpression{
				Left: expression.OrExpression{
					Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(0, 2)}},
					Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(7, 9)}},
				},
				Right: expression.TokenStream{expression.Token{Value: "three", Span: span.NewSpan(14, 18)}},
			},
		},
		{
			name: "mixed operators are left-associative with no precedence",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				l("two", span.NewSpan(8, 10)),
				or(span.NewSpan(12, 13)),
				l("three", span.NewSpan(15, 19)),
			},
			expected: expression.OrExpression{
				Left: expression.AndExpression{
					Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(0, 2)}},
					Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(8, 10)}},
				},
				Right: expression.TokenStream{expression.Token{Value: "three", Span: span.NewSpan(15, 19)}},
			},
		},
		{
			name: "not as right operand of and",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				not(span.NewSpan(8, 10)),
				l("two", span.NewSpan(12, 14)),
			},
			expected: expression.AndExpression{
				Left: expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(0, 2)}},
				Right: expression.NotExpression{
					Expr: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(12, 14)}},
				},
			},
		},
		{
			name: "not as right operand of or",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				or(span.NewSpan(4, 5)),
				not(span.NewSpan(7, 9)),
				l("two", span.NewSpan(11, 13)),
			},
			expected: expression.OrExpression{
				Left: expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(0, 2)}},
				Right: expression.NotExpression{
					Expr: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(11, 13)}},
				},
			},
		},
		{
			name: "not with parenthesized group as right operand",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
				and(span.NewSpan(4, 6)),
				not(span.NewSpan(8, 10)),
				lp(span.NewSpan(12, 12)),
				l("two", span.NewSpan(13, 15)),
				or(span.NewSpan(17, 18)),
				l("three", span.NewSpan(20, 24)),
				rp(span.NewSpan(25, 25)),
			},
			expected: expression.AndExpression{
				Left: expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(0, 2)}},
				Right: expression.NotExpression{
					Expr: expression.OrExpression{
						Left:  expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(13, 15)}},
						Right: expression.TokenStream{expression.Token{Value: "three", Span: span.NewSpan(20, 24)}},
					},
				},
			},
		},
		{
			name:        "empty input",
			input:       lexer.TokenStream{},
			expectedErr: newParserError(ErrInvalidSyntax, span.NewSpan(0, 0)),
		},
		{
			name: "single literal",
			input: lexer.TokenStream{
				l("one", span.NewSpan(0, 2)),
			},
			expectedErr: newParserError(fmt.Errorf("%w: %w", ErrInvalidSyntax, ErrToExpressionUnknownNodeKind), span.NewSpan(0, 2)),
		},

		{
			name: "multiword",
			input: lexer.TokenStream{
				q(span.NewSpan(0, 0)),
				l("one", span.NewSpan(1, 3)),
				l("two", span.NewSpan(5, 7)),
				q(span.NewSpan(8, 8)),
				and(span.NewSpan(10, 12)),
				l("three", span.NewSpan(14, 18)),
			},
			expected: expression.AndExpression{
				Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(1, 3)}, expression.Token{Value: "two", Span: span.NewSpan(5, 7)}},
				Right: expression.TokenStream{expression.Token{Value: "three", Span: span.NewSpan(14, 18)}},
			},
		},
		{
			name: "bare parenthesized expression at root",
			input: lexer.TokenStream{
				lp(span.NewSpan(0, 0)),
				l("one", span.NewSpan(1, 3)),
				and(span.NewSpan(5, 7)),
				l("two", span.NewSpan(9, 11)),
				rp(span.NewSpan(12, 12)),
			},
			expected: expression.AndExpression{
				Left:  expression.TokenStream{expression.Token{Value: "one", Span: span.NewSpan(1, 3)}},
				Right: expression.TokenStream{expression.Token{Value: "two", Span: span.NewSpan(9, 11)}},
			},
		},
		{
			name: "not followed by not",
			input: lexer.TokenStream{
				not(span.NewSpan(0, 2)),
				not(span.NewSpan(4, 6)),
				l("one", span.NewSpan(8, 10)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxNotInvalidFollowing, span.NewSpan(4, 6)),
		},
		{
			name: "empty parentheses",
			input: lexer.TokenStream{
				lp(span.NewSpan(0, 0)),
				rp(span.NewSpan(1, 1)),
			},
			expectedErr: newParserError(fmt.Errorf("%w: %w", ErrInvalidSyntax, ErrToExpressionUnknownNodeKind), span.NewSpan(0, 1)),
		},
		{
			name: "standalone and keyword",
			input: lexer.TokenStream{
				and(span.NewSpan(0, 2)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxAndNoLeftExpr, span.NewSpan(0, 2)),
		},
		{
			name: "standalone or keyword",
			input: lexer.TokenStream{
				or(span.NewSpan(0, 1)),
			},
			expectedErr: newParserError(ErrInvalidSyntaxOrNoLeftExpr, span.NewSpan(0, 1)),
		},
	}

	c := parser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := c.Parse(tt.input)
			if tt.expectedErr == nil {
				require.NoError(t, err)
				require.Equal(t, tt.expected, expr)
				return
			}

			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}
