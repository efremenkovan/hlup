package lexer

import "github.com/efremenkovan/hlup/span"

type TokenKind int8

const (
	TokenKindLPar       TokenKind = iota // (
	TokenKindRPar                        // )
	TokenKindQuote                       // ' or "
	TokenKindKeywordOR                   // "or"
	TokenKindKeywordAND                  // "and"
	TokenKindKeywordNOT                  // "not"
	TokenKindLiteral                     // any string/number besides declared above
)

type TokenStream []Token

// Token representes a single meaningfull piece of expression
type Token struct {
	kind  TokenKind
	value string
	span  span.Span
}

func NewToken(kind TokenKind, span span.Span) Token {
	return Token{
		kind:  kind,
		value: "",
		span:  span,
	}
}

// WithSpan mutates token, attaching a span details to it
func (t *Token) WithSpan(span span.Span) {
	t.span = span
}

func NewLiteralToken(value string, span span.Span) Token {
	return Token{
		kind:  TokenKindLiteral,
		value: value,
		span:  span,
	}
}

// String returns a string representation of token
func (t *Token) String() string {
	switch t.kind {
	case TokenKindLPar:
		return "("
	case TokenKindRPar:
		return ")"
	case TokenKindQuote:
		return "\""
	case TokenKindKeywordAND:
		return "AND"
	case TokenKindKeywordNOT:
		return "NOT"
	case TokenKindKeywordOR:
		return "OR"
	case TokenKindLiteral:
		return t.value
	}

	return "<UNKNOWN TOKEN>"
}

func (t Token) Kind() TokenKind {
	return t.kind
}

func (t Token) Span() span.Span {
	return t.span
}

// Matches an input string to respecteful token
// If no known string value is found, treat input as Literal
func tokenFromString(set *keywordSet, input string, span span.Span) Token {
	if v, ok := set.KeywordTokenKind(input); ok {
		return NewToken(v, span)
	}

	switch input {
	case "(":
		return NewToken(TokenKindLPar, span)
	case ")":
		return NewToken(TokenKindRPar, span)
	case "\"", "'":
		return NewToken(TokenKindQuote, span)
	}

	return NewLiteralToken(input, span)
}
