package benchmark

import (
	"fmt"
	"testing"

	"github.com/efremenkovan/hlup"
	"github.com/efremenkovan/hlup/expression"
	"github.com/efremenkovan/hlup/internal/benchmark/testdata"
	"github.com/efremenkovan/hlup/lexer"
	"github.com/efremenkovan/hlup/parser"
)

func BenchmarkTokenizeInput(b *testing.B) {
	type benchmark struct {
		name  string
		input string
	}

	benchmarks := []benchmark{
		{
			name:  "short input",
			input: testdata.TextShort,
		},
		{
			name:  "medium input",
			input: testdata.TextMedium,
		},
		{
			name:  "long input",
			input: testdata.TextLong,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()

			for b.Loop() {
				hlup.TokenizeInput(bm.input)
			}
		})
	}
}

func BenchmarkLexExpression(b *testing.B) {
	type benchmark struct {
		name  string
		input string
	}

	benchmarks := []benchmark{
		{
			name:  "simple query",
			input: testdata.QuerySimple,
		},
		{
			name:  "medium query",
			input: testdata.QueryMedium,
		},
		{
			name:  "complex query",
			input: testdata.QueryComplex,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			lexer := lexer.NewLexer(bm.input)

			b.ResetTimer()

			for b.Loop() {
				lexer.TokenStream()
			}
		})
	}
}

func BenchmarkCompileLexedExpression(b *testing.B) {
	type benchmark struct {
		name  string
		input lexer.TokenStream
	}

	benchmarks := []benchmark{
		{
			name:  "simple query",
			input: lexer.NewLexer(testdata.QuerySimple).TokenStream(),
		},
		{
			name:  "medium query",
			input: lexer.NewLexer(testdata.QueryMedium).TokenStream(),
		},
		{
			name:  "complex query",
			input: lexer.NewLexer(testdata.QueryComplex).TokenStream(),
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			parser := parser.NewParser()

			b.ResetTimer()

			for b.Loop() {
				_, _ = parser.Parse(bm.input)
			}
		})
	}
}

func BenchmarkCompileExpressionString(b *testing.B) {
	type benchmark struct {
		name  string
		input string
	}

	benchmarks := []benchmark{
		{
			name:  "simple query",
			input: testdata.QuerySimple,
		},
		{
			name:  "medium query",
			input: testdata.QueryMedium,
		},
		{
			name:  "complex query",
			input: testdata.QueryComplex,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()

			for b.Loop() {
				_, _ = hlup.CompileExpression(bm.input)
			}
		})
	}
}

func BenchmarkMatch(b *testing.B) {
	type benchmark struct {
		name           string
		inputTokenized expression.TokenStream
		expression     expression.Expression
	}

	shortTokenizedText := hlup.TokenizeInput(testdata.TextShort)
	mediumTokenizedText := hlup.TokenizeInput(testdata.TextMedium)
	longTokenizedText := hlup.TokenizeInput(testdata.TextLong)

	simpleExpression, _ := hlup.CompileExpression(testdata.QuerySimple)
	mediumExpression, _ := hlup.CompileExpression(testdata.QueryMedium)
	complexExpression, _ := hlup.CompileExpression(testdata.QueryComplex)

	texts := map[string]expression.TokenStream{
		"short":  shortTokenizedText,
		"medium": mediumTokenizedText,
		"long":   longTokenizedText,
	}
	expressions := map[string]expression.Expression{
		"simple":  simpleExpression,
		"medium":  mediumExpression,
		"complex": complexExpression,
	}

	benchmarks := make([]benchmark, 0, len(texts)*len(expressions))

	for tname, text := range texts {
		for ename, expr := range expressions {
			benchmarks = append(benchmarks, benchmark{
				name:           fmt.Sprintf("%s text, %s expression", tname, ename),
				expression:     expr,
				inputTokenized: text,
			})
		}
	}

	b.Run("match contains", func(b *testing.B) {
		for _, bm := range benchmarks {
			b.Run(bm.name, func(b *testing.B) {
				b.ResetTimer()

				for b.Loop() {
					bm.expression.MatchContains(bm.inputTokenized)
				}
			})
		}
	})

	b.Run("match equal", func(b *testing.B) {
		for _, bm := range benchmarks {
			b.Run(bm.name, func(b *testing.B) {
				b.ResetTimer()

				for b.Loop() {
					bm.expression.MatchEqual(bm.inputTokenized)
				}
			})
		}
	})
}
