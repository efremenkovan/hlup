# Hlup

Hlup lets you search for word combinations (queries) in text.
It's small, extendable, fast and have zero dependencies.


[![go.mod](https://img.shields.io/github/go-mod/go-version/efremenkovan/hlup)](go.mod)
[![Go Reference](https://pkg.go.dev/badge/github.com/efremenkovan/hlup.svg)](https://pkg.go.dev/github.com/efremenkovan/hlup)
[![Go Report Card](https://goreportcard.com/badge/github.com/efremenkovan/hlup)](https://goreportcard.com/report/github.com/efremenkovan/hlup)
[![CI Lint](https://github.com/efremenkovan/hlup/actions/workflows/lint.yml/badge.svg)](https://github.com/efremenkovan/hlup/actions/workflows/lint.yml)
[![CI Test](https://github.com/efremenkovan/hlup/actions/workflows/test.yml/badge.svg)](https://github.com/efremenkovan/hlup/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/license-MIT-orange.svg?style=flat)](https://github.com/efremenkovan/hlup/blob/main/LICENSE)

## Navigation

- [Features](#features)
- [Installation](#installation)
- [Query syntax](#query-syntax)
  - [Languages support](#languages-support)
- [Example usage](#example-usage)
- [Matching strictness](#matching-strictness)
  - [Equality check](#equality-check)
  - [Contains check](#contains-check)
- [Text highlight with spans](#text-highlight-with-spans)
- [Benchmarks](#benchmarks)
- [Input text tokenization](#input-text-tokenization)
  - [Word breaker symbols](#word-breaker-symbols)
  - [Accent stripping](#accent-stripping)

## Features

- Simple query language
- Highlight matched text pieces
- Simple caching friendly API
- Reasonable defaults with extension possibilities
- No external dependencies
- [It's fast](#benchmarks)

## Installation

```bash
go get github.com/efremenkovan/hlup
```

## Query syntax

Hlup supports primitive query language with 3 expressions:

- and - `one and two`
- or  - `one or two`
- not - `not one`

- They can be chained together:
`one and not two or three`

- They can be nested and prioritized via parentheses:
`(one and two) or (three and four)`

- Whole phrase can be declared via quotes (both single and double):
`("first you have to read" and write) and code`. This way hlup will look for exactly
"first you have to read" phrase in text.

### Languages support

As of now hlup supports 2 languages in queries:

- English - and/or/not
- Russian - и/или/не

You can specify a language you wrote your query in via `WithLang` func
from `github.com/efremenkovan/hlup/options` package when
calling `CompileExpression` function:

```go
import (
    "github.com/efremenkovan/hlup/options"
    "github.com/efremenkovan/hlup/lang"
)
hlup.CompileExpression("раз и два", options.WithLang(lang.LangRU))
```

English is the default language, no need to specify it with options.

> [!IMPORTANT]
> It only affects keywords. You can use any language you wish for the words you
> want to look for, but `and/or/not` logical expressions can only be written in
> one of the two supported languages.
>
> It means you can write `geschwindigkeit and begrenzung`, using German words and
> english logic expressions.

## Example usage

```go
expression := hlup.CompileExpression("go and programming and not 'to hell'")
tokenizedText := hlup.TokenizeInput("programming in go")

matches, spans := expression.MatchEqual(tokenizedText)

println(matches) // -> true
fmt.Printf("%v\n", spans) // -> []span.Span{ {Start: 15, End: 16}, {Start: 0, End: 10} }
```

> [!TIP]
> Always prefer to compile each query and tokenize input text exactly once and store
> compiled expressions and tokenized texts in some kind of cache.
> Despite hlup compilation and tokenization being fast, they require computation.

## Matching strictness

There are 2 ways hlup matches tokens in queries:

- equality check
- contains check

### Equality check

Each expression token should be equal to respective input token:

Given a query `make and bam`

It will match `we make bam sometimes`, since exactly `make` and `bam`
are present in text.

It will not match `we make bamboo sometimes`, since `make` is present in text,
but `bam` is not.
`bamboo` != `bam`

### Contains check

Each expression token should be included in respective input token:
Given a query `make and bam`

It will match `we make bam sometimes`, since exactly `make` and `bam`
are present in text.

It will also match `we make bamboo sometimes`, since `make` is present in text,
and `bamboo` contains `bam` in it

## Text highlight with spans

If query matches the provided text, match function returns a slice of `Span` objects.

```go
type Span struct {
    Start int // index of the first matched character in original text
    End   int // index of the last matched character in original text
}
```

You can use them to highlight matched words.

> [!NOTE]
> Span's Start and End fields are INDEXES of characters to look for in original text
> For example "I" in "I give you my word" will have a span of `{Start: 0, End: 0}`
> and "give" will have a span of `{Start: 2, End: 5}`

> [!IMPORTANT]
> Span slice is sorted by query word sequence, not the original text word sequence.

Example:
Given a query `go and programming` matching against `programming in go` will give
you following slice of spans:
`[]span.Span{ {Start: 15, End: 16}, {Start: 0, End: 10} }`

If you would want to highlight the whole `programming in go` word sequence, you would
have to find the most left and the most right span positions yourself.

## Benchmarks

All the measurements are done on Macbook Pro 16" Apple M3 Pro

### Text tokenization

|Text|Time|
|-----------------|-----|
|short (~30 words)|3.6µs|
|medium (~1k words)|103µs|
|long (~10k words)|1.19ms|

### Query compilation

|Query|Time|
|-----------------|-----|
|simple (4 words, no nesting)|832ns|
|medium (~20 words, 4 lvl nesting)|5.7µs|
|complex (~40 words, 12 lvl nesting)|17.5µs|

### Query matching

#### Match contains

|Expr|Short text|Medium text|Long text|
|----|----------|-----------|---------|
|simple|375ns|6.3µs|86µs|
|medium|241ns|6.2µs|89µs|
|complex|446ns|10.4µs|157µs|

#### Match equal

|Expr|Short text|Medium text|Long text|
|----|----------|-----------|---------|
|simple|196ns|2.85µs|28µs|
|medium|96ns|2.72µs|28µs|
|complex|205ns|5.29µs|56µs|

## Input text tokenization

Input text tokenization means splitting the plain text to separate words to be
matched against the query.

For hlup to be more effective in terms of giving reasonable match results, it
transforms text to lowercase, splits it into separate words and performs
accent stripping.

For English and Russian languages hlup has reasonable defaults, you don't have
to do anything to get good match results.

But if you want to match texts written in any other language, you can extend hlup
tokenization settings for that.

### Word breaker symbols

To be able to perform text tokenization we have to distinguish some word boundaries.
There is a bunch of characters considered a word breakers by default:

- Punctuation - `! ? : ; - — . , * { } ( ) [ ] +`
- Quotes      - `` ` « » " ' ``
- Slashes     - `/ | \`
- Whitespaces - `\t \n \v \f \r   0x85  0xA0`

> [!IMPORTANT]
> All word breaker characters are dropped from tokenized text. You can't match
> against them

> [!TIP]
> You can extend list of word breaker characters:

```go
hlup.TokenizeInput(
    textYouWannaCheck,
    hlup.WithExtendedWordBreakersList([]rune{ '%', '^' }),
)
```

You can also completely override list of word breaker characters via
`hlup.WithCustomWordBreakersList`

### Accent stripping

There are some characters that you might want to treat like others.

```
ö -> o
ä -> a
ü -> u
ô -> o
```

... and so on


Hlup performs accent stripping during tokenization process.
Since hlup supports only English and Russian out of the box, replacement table
is as simple as that:

```go
map[rune]rune{
    'ё': 'е',
}
```

> [!TIP]
> You can extend replacement characters table via custom options:

```go
hlup.TokenizeInput(
    textYouWannaCheck,
    hlup.WithExtendedReplaceTable(map[rune]rune{
        'ö': 'o',
        'ä': 'a',
        'ü': 'u',
        'ô': 'o',
    }),
)
```

That way you can improve matching for any language you want.

You can also completely override the replacement table via `hlup.WithCustomReplaceTable`.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
