package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/efremenkovan/hlup/lang"
	"github.com/efremenkovan/hlup/options"
	"github.com/efremenkovan/hlup/span"
)

type lexer struct {
	kSet  *keywordSet
	input []rune
	cur   cursor
	ctx   ctx
}

type cursor struct {
	position int
}

// SlideByOne mutates cursor state by moving it 1 position further
func (c *cursor) SlideByOne() {
	c.position += 1
}

// Clone creates a new instance of cursor with state identical to the referenced one
func (c cursor) Clone() cursor {
	return cursor{
		position: c.position,
	}
}

// Returns new cursor at 0 position
func newCursor() cursor {
	return cursor{
		position: 0,
	}
}

type Lexer interface {
	// TokenStream returns a list of lexer tokens representing whole input string
	// TokenStream does not adjust cursor position, it can be called at any time.
	TokenStream() TokenStream
	// Consume returns next token and adjusts cursor position
	Consume() (token Token, ok bool)
	// Peek returns next token without adjusting cursor position
	Peek() (token Token, ok bool)
	// IsFinished returns true if there is no more tokens to consume
	IsFinished() bool
}

// NewLexer returns new instance of lexer
// Lexer lang is english by default. To change that use WithLang option
func NewLexer(input string, opts ...options.PatchFunc) *lexer {
	o := options.Options{
		Lang: lang.LangEN,
	}

	for _, f := range opts {
		f(&o)
	}

	return &lexer{
		input: []rune(input),
		kSet:  getSetByLang(o.Lang),
		ctx:   newCtx(),
		cur:   newCursor(),
	}
}

// Consume returns next token and adjusts cursor position
func (l *lexer) Consume() (Token, bool) {
	v, ok := l.takeTokenAndAdvanceCursor(&l.ctx, &l.cur)
	l.SkipWhitespaces()
	return v, ok
}

// Adjusts cursor position to be placed at first found non-whitespace rune
func (l *lexer) SkipWhitespaces() {
	l.skipWhitespacesOnCursor(&l.cur)
}

// Adjusts cursor position to be placed at first found non-whitespace rune
func (l *lexer) skipWhitespacesOnCursor(cur *cursor) {
	for cur.position < len(l.input) {
		ch := l.input[cur.position]
		if unicode.IsSpace(ch) {
			cur.SlideByOne()
		} else {
			break
		}
	}
}

// Peek returns next token without adjusting cursor position
func (l *lexer) Peek() (Token, bool) {
	// create copy of a context to not mutate lexer real context
	ctx := l.ctx.Clone()
	// create copy of a cursor to not advance lexer cursor position
	cur := l.cur.Clone()

	return l.takeTokenAndAdvanceCursor(&ctx, &cur)
}

func (l *lexer) TokenStream() TokenStream {
	stream := make(TokenStream, 0)
	ctx := newCtx()
	cur := newCursor()

	for cur.position < len(l.input) {
		if token, ok := l.takeTokenAndAdvanceCursor(&ctx, &cur); ok {
			stream = append(stream, token)
			l.skipWhitespacesOnCursor(&cur)
		} else {
			break
		}
	}

	return stream
}

// Consumes following token, advancing cursor position and mutating the context as it goes
func (l *lexer) takeTokenAndAdvanceCursor(
	ctx *ctx,
	cur *cursor,
) (Token, bool) {
	buffer := make([]rune, 0)

loop:
	for cur.position < len(l.input) {
		ch := l.input[cur.position]

		// no matter what char it is, append it to buffer since it's escaped
		if ctx.Is(ctxEscapedRune) {
			buffer = append(buffer, ch)
			cur.SlideByOne()
			// only one char is escaped by reverse slash
			ctx.Drop(ctxEscapedRune)
			continue
		}

		switch ch {
		case '\\':
			ctx.Add(ctxEscapedRune)
			cur.SlideByOne()
			continue
		case ')', '(':
			// parentheses are terminal characters for literal tokens
			if len(buffer) > 0 {
				break loop
			}

			buffer = append(buffer, ch)
			cur.SlideByOne()
			break loop
		case '\'', '"':
			// Quote acts like a terminal character for multiword literal tokens
			if ctx.Is(ctxInQuote) && ch == ctx.terminalQuote && len(buffer) > 0 {
				break loop
			}

			buffer = append(buffer, ch)
			cur.SlideByOne()

			// Leave multiword literal token context
			if ctx.Is(ctxInQuote) {
				ctx.terminalQuote = rune(0)
				ctx.Drop(ctxInQuote)
				break loop
			}

			// Enter multiword literal token context
			ctx.terminalQuote = ch
			ctx.Add(ctxInQuote)
			break loop
		case ' ':
			break loop
		default:
			buffer = append(buffer, ch)
			cur.SlideByOne()
			continue
		}
	}

	if len(buffer) == 0 {
		return Token{}, false
	}

	str := string(buffer)
	token := tokenFromString(l.kSet, str, span.NewSpan(cur.position-utf8.RuneCountInString(str), cur.position-1))

	return token, true
}

// IsFinished returns true if there is no more tokens to consume
func (l *lexer) IsFinished() bool {
	return l.cur.position >= len(l.input)
}
