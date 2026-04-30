package hlup

import (
	"maps"
	"strings"

	"github.com/efremenkovan/hlup/expression"
	"github.com/efremenkovan/hlup/lexer"
	"github.com/efremenkovan/hlup/options"
	"github.com/efremenkovan/hlup/parser"
	"github.com/efremenkovan/hlup/span"
)

// CompileExpression returns an expression that represents provided rule in specified language.
// NOTE: Default language is english. To change the language use options.WithLang
func CompileExpression(rule string, optionsPatches ...options.PatchFunc) (expression.Expression, error) {
	lexer := lexer.NewLexer(rule, optionsPatches...)
	parser := parser.NewParser()
	return parser.Parse(lexer.TokenStream())
}

var defaultReplaceTable = map[rune]rune{
	'ё': 'е',
}

var defaultWordBreakersTable = map[rune]struct{}{
	'!': {},
	'?': {},
	':': {},
	';': {},
	'-': {},
	'—': {},
	'.': {},
	',': {},
	'*': {},
	'+': {},
	'(': {},
	')': {},
	'{': {},
	'}': {},
	'[': {},
	']': {},

	// Quotes
	'«':  {},
	'»':  {},
	'"':  {},
	'\'': {},
	'`':  {},

	// Slashes
	'/':  {},
	'\\': {},
	'|':  {},

	// Whitespaces
	'\t': {},
	'\n': {},
	'\v': {},
	'\f': {},
	'\r': {},
	' ':  {},
	0x85: {},
	0xA0: {},
}

type tokenizeOptions struct {
	replaceTable      map[rune]rune
	wordBreakersTable map[rune]struct{}
}

type tokenizeOptionsPatchFunc func(o *tokenizeOptions)

// WithExtendedReplaceTable extends rune replace lookup table with provided onces.
func WithExtendedReplaceTable(table map[rune]rune) tokenizeOptionsPatchFunc {
	return func(o *tokenizeOptions) {
		maps.Copy(o.replaceTable, table)
	}
}

// WithCustomReplaceTable overrides rune replace lookup table with provided onces.
// Using this function is not recommended, consider using WithExtendedReplaceTable.
func WithCustomReplaceTable(table map[rune]rune) tokenizeOptionsPatchFunc {
	return func(o *tokenizeOptions) {
		o.replaceTable = table
	}
}

// WithExtendedWordBreakersList extends word breaker runes set with list of provided ones.
func WithExtendedWordBreakersList(list []rune) tokenizeOptionsPatchFunc {
	return func(o *tokenizeOptions) {
		for _, rune := range list {
			o.wordBreakersTable[rune] = struct{}{}
		}
	}
}

// WithCustomWordBreakersList overrides all known word breaker runes.
// Using this function is not recommended, consider using WithExtendedWordBreakersList.
func WithCustomWordBreakersList(list []rune) tokenizeOptionsPatchFunc {
	return func(o *tokenizeOptions) {
		o.wordBreakersTable = make(map[rune]struct{}, len(list))
		for _, rune := range list {
			o.wordBreakersTable[rune] = struct{}{}
		}
	}
}

// TokenizeInput returns a token (word) sequence, dropping punctuation.
func TokenizeInput(input string, options ...tokenizeOptionsPatchFunc) expression.TokenStream {
	opts := tokenizeOptions{
		wordBreakersTable: defaultWordBreakersTable,
		replaceTable:      defaultReplaceTable,
	}

	for _, patchFunc := range options {
		patchFunc(&opts)
	}

	input = strings.ToLower(input)

	runes := []rune(input)
	buffer := make([]rune, 0, 12)
	bufferStart := -1
	bufferEnd := -1

	result := make(expression.TokenStream, 0)

	for index, r := range runes {
		if _, ok := opts.wordBreakersTable[r]; ok {
			// Persist token
			if len(buffer) > 0 && bufferStart >= 0 {
				result = append(result, expression.Token{
					Value: string(buffer),
					Span:  span.NewSpan(bufferStart, bufferEnd),
				})
			}

			// Reset buffer state
			buffer = buffer[:0]
			bufferStart = -1
			bufferEnd = -1
			continue
		}

		if newR, ok := opts.replaceTable[r]; ok {
			r = newR
		}

		if bufferStart == -1 {
			bufferStart = index
		}

		bufferEnd = index

		buffer = append(buffer, r)
	}

	// We exited loop due to EOF but there was no terminal rune to put buffer into results
	if len(buffer) > 0 {
		result = append(result, expression.Token{
			Value: string(buffer),
			Span:  span.NewSpan(bufferStart, bufferEnd),
		})
	}

	return result
}
