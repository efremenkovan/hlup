package parser

import (
	"errors"
	"fmt"

	"github.com/efremenkovan/hlup/span"
)

type ParserError struct {
	cause error
	span  span.Span
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s near: %d..%d", e.cause, e.span.Start, e.span.End)
}

func (e *ParserError) Unwrap() error {
	return e.cause
}

func newParserError(msg error, span span.Span) *ParserError {
	return &ParserError{
		cause: msg,
		span:  span,
	}
}

type SyntaxError struct {
	message string
	cause   error
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", e.cause, e.message)
}

func (e *SyntaxError) Unwrap() error {
	return e.cause
}

var (
	ErrInsertChildToNode = errors.New("failed to insert child to node")

	ErrToExpressionNoLeftNode      = errors.New("no left node")
	ErrToExpressionNoRightNode     = errors.New("no right node")
	ErrToExpressionUnknownNodeKind = errors.New("unknown node kind")

	ErrInvalidSyntax                      = errors.New("invalid syntax")
	ErrInvalidSyntaxUnbalancedParentheses = &SyntaxError{message: "unbalanced parentheses", cause: ErrInvalidSyntax}
	ErrInvalidSyntaxNotNoFollowingExpr    = &SyntaxError{message: "NOT expression have no following expression to negate", cause: ErrInvalidSyntax}
	ErrInvalidSyntaxNotInvalidFollowing   = &SyntaxError{message: "NOT expression is followed by keyword", cause: ErrInvalidSyntax}
	ErrInvalidSyntaxAndNoLeftExpr         = &SyntaxError{message: "AND expression have no left branch", cause: ErrInvalidSyntax}
	ErrInvalidSyntaxAndNoRightExpr        = &SyntaxError{message: "AND expression have no right branch", cause: ErrInvalidSyntax}
	ErrInvalidSyntaxAndInvalidFollowing   = &SyntaxError{message: "AND expression is is followed by unsupported keyword", cause: ErrInvalidSyntax}
	ErrInvalidSyntaxOrNoLeftExpr          = &SyntaxError{message: "OR expression have no left branch", cause: ErrInvalidSyntax}
	ErrInvalidSyntaxOrNoRightExpr         = &SyntaxError{message: "OR expression have no right branch", cause: ErrInvalidSyntax}
	ErrInvalidSyntaxOrInvalidFollowing    = &SyntaxError{message: "OR expression is is followed by unsupported keyword", cause: ErrInvalidSyntax}
)

var (
	ErrTreeInsertIntoLeaf     = errors.New("leaf can have no child nodes")
	ErrTreeInsertIntoFullNode = errors.New("node is at child nodes limit")
)
