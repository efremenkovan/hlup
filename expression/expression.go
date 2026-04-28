package expression

import (
	"slices"

	"github.com/efremenkovan/hlup/span"
)

type Expression interface {
	// MatchEqual returns true if the input contains contiguous subsequence, where each token is equal to the corresponding token in the expression
	MatchEqual(input TokenStream) (found bool, spans []span.Span)

	// MatchEqual returns true if the input contains contiguous subsequence, where each token contains the corresponding token in the expression
	MatchContains(input TokenStream) (found bool, spans []span.Span)
}

type matchFunc func(input TokenStream) (found bool, spans []span.Span)

// AndExpression is true only when both - Left and Right expression branches are true
type AndExpression struct {
	Left  Expression
	Right Expression
}

func (e AndExpression) MatchContains(tokenSequence TokenStream) (bool, []span.Span) {
	return e.match(tokenSequence, e.Left.MatchContains, e.Right.MatchContains)
}

func (e AndExpression) MatchEqual(tokenSequence TokenStream) (bool, []span.Span) {
	return e.match(tokenSequence, e.Left.MatchEqual, e.Right.MatchEqual)
}

func (e AndExpression) match(
	tokenSequence TokenStream,
	lFunc, rFunc matchFunc,
) (bool, []span.Span) {
	if isLeftMatch, leftSpan := lFunc(tokenSequence); isLeftMatch {
		if isRightMatch, rightSpan := rFunc(tokenSequence); isRightMatch {
			return true, slices.Concat(leftSpan, rightSpan)
		}
	}

	return false, nil
}

// OrExpression is true when either the Left or the Right expression branch is true
type OrExpression struct {
	Left  Expression
	Right Expression
}

func (e OrExpression) MatchContains(tokenSequence TokenStream) (bool, []span.Span) {
	return e.match(tokenSequence, e.Left.MatchContains, e.Right.MatchContains)
}

func (e OrExpression) MatchEqual(tokenSequence TokenStream) (bool, []span.Span) {
	return e.match(tokenSequence, e.Left.MatchEqual, e.Right.MatchEqual)
}

func (e OrExpression) match(
	tokenSequence TokenStream,
	lFunc, rFunc matchFunc,
) (bool, []span.Span) {
	if isLeftMatch, leftSpan := lFunc(tokenSequence); isLeftMatch {
		return true, leftSpan
	}

	if isRightMatch, rightSpan := rFunc(tokenSequence); isRightMatch {
		return true, rightSpan
	}

	return false, nil
}

// NotExpression is only true when child expression is false
type NotExpression struct {
	Expr Expression
}

func (e NotExpression) MatchContains(tokenSequence TokenStream) (bool, []span.Span) {
	return e.match(tokenSequence, e.Expr.MatchContains)
}

func (e NotExpression) MatchEqual(tokenSequence TokenStream) (bool, []span.Span) {
	return e.match(tokenSequence, e.Expr.MatchEqual)
}

func (e NotExpression) match(tokenSequence TokenStream, f matchFunc) (bool, []span.Span) {
	if match, value := f(tokenSequence); match {
		return false, value
	}

	return true, nil
}
