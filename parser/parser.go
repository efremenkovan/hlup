package parser

import (
	"fmt"

	"github.com/efremenkovan/hlup/expression"
	"github.com/efremenkovan/hlup/lexer"
	"github.com/efremenkovan/hlup/span"
)

type parser struct{}

// Parser is used to transform lexer output to expression tree to be matched against
type Parser interface {
	// Parse returns expression tree, represeting input token sequence
	Parse(input lexer.TokenStream) (expression.Expression, error)
}

func NewParser() Parser {
	return &parser{}
}

// Parse returns expression tree, represeting input token sequence
func (p *parser) Parse(input lexer.TokenStream) (expression.Expression, error) {
	if len(input) == 0 {
		return nil, newParserError(ErrInvalidSyntax, span.NewSpan(0, 0))
	}
	rootNode := newNode()
	currentNode := rootNode
	currentNestLevel := 0

	for tokenIdx, token := range input {
		if currentNode == nil {
			break
		}
		isLastToken := tokenIdx == len(input)-1

		switch token.Kind() {
		case lexer.TokenKindQuote:
			if currentNode.IsFull() {
				return nil, newParserError(ErrInvalidSyntax, token.Span())
			}

			// We are already inside quoted literal sequence, just return
			if currentNode.Kind() == nodeKindLeaf {
				currentNode = currentNode.Parent()

				// We have to exit not node when we inserted a token into it
				if currentNode.Kind() == nodeKindNOT {
					currentNode = currentNode.Parent()
				}
				continue
			}

			leaf := emptyLeaf()
			leaf.WithParent(currentNode)
			if err := currentNode.Insert(&leaf); err != nil {
				return nil, newParserError(err, token.Span())
			}
			currentNode = &leaf

		case lexer.TokenKindLPar:
			node := newNode()
			node.WithSpan(token.Span())
			node.WithParent(currentNode)
			if err := currentNode.Insert(node); err != nil {
				return nil, newParserError(fmt.Errorf("%w: %w", ErrInsertChildToNode, err), token.Span())
			}
			currentNestLevel += 1
			currentNode = node

		case lexer.TokenKindRPar:
			if currentNode.IsRoot() {
				return nil, newParserError(ErrInvalidSyntaxUnbalancedParentheses, token.Span())
			}

			currentNode.WithSpan(span.NewSpan(currentNode.Span().Start, token.Span().End))
			currentNestLevel -= 1

			// We have to exit not node when we inserted a token into it
			if currentNode.Kind() == nodeKindNOT {
				currentNode = currentNode.Parent()
			}
			currentNode = currentNode.Parent()

		case lexer.TokenKindLiteral:
			if currentNode.IsFull() {
				return nil, newParserError(ErrInvalidSyntax, token.Span())
			}

			// We are inside quoted literal sequence
			if currentNode.Kind() == nodeKindLeaf {
				if leafNode, ok := currentNode.(*leaf); ok {
					leafNode.AppendContent(token)
				} else {
					return nil, newParserError(ErrInvalidSyntax, token.Span())
				}
				continue
			}

			leaf := newLeaf(token)
			if err := currentNode.Insert(&leaf); err != nil {
				return nil, newParserError(err, token.Span())
			}

			// We have to exit not node when we inserted a token into it
			if currentNode.Kind() == nodeKindNOT {
				currentNode = currentNode.Parent()
			}

		case lexer.TokenKindKeywordNOT:
			if isLastToken {
				return nil, newParserError(ErrInvalidSyntaxNotNoFollowingExpr, token.Span())
			}

			followingToken := input[tokenIdx+1]
			switch followingToken.Kind() {
			case lexer.TokenKindKeywordAND, lexer.TokenKindKeywordNOT, lexer.TokenKindKeywordOR:
				return nil, newParserError(ErrInvalidSyntaxNotInvalidFollowing, followingToken.Span())
			}

			if currentNode.IsRoot() && currentNode.IsEmpty() {
				currentNode.WithKind(nodeKindNOT)
				currentNode.WithSpan(token.Span())
				continue
			}

			node := newNode()
			node.WithKind(nodeKindNOT)
			node.WithParent(currentNode)
			node.WithSpan(token.Span())
			if err := currentNode.Insert(node); err != nil {
				return nil, newParserError(fmt.Errorf("%w: %w", ErrInsertChildToNode, err), token.Span())
			}
			currentNode = node

		case lexer.TokenKindKeywordOR:
			if currentNode.Kind() == nodeKindUnknown {
				if currentNode.IsEmpty() {
					return nil, newParserError(ErrInvalidSyntaxOrNoLeftExpr, token.Span())
				}

				currentNode.WithKind(nodeKindOR)
				currentNode.WithSpan(token.Span())

				if isLastToken {
					return nil, newParserError(ErrInvalidSyntaxOrNoRightExpr, token.Span())
				}

				followingToken := input[tokenIdx+1]
				switch followingToken.Kind() {
				case lexer.TokenKindKeywordAND, lexer.TokenKindKeywordOR:
					return nil, newParserError(ErrInvalidSyntaxOrInvalidFollowing, followingToken.Span())
				}

				continue
			}

			node := currentNode.CloneDetached()
			node.WithParent(currentNode)
			node.WithSpan(token.Span())

			currentNode.ClearContent()
			currentNode.WithKind(nodeKindOR)
			if err := currentNode.Insert(node); err != nil {
				return nil, newParserError(fmt.Errorf("%w: %w", ErrInsertChildToNode, err), token.Span())
			}

		case lexer.TokenKindKeywordAND:
			if currentNode.Kind() == nodeKindUnknown {
				if currentNode.IsEmpty() {
					return nil, newParserError(ErrInvalidSyntaxAndNoLeftExpr, token.Span())
				}
				currentNode.WithKind(nodeKindAND)
				currentNode.WithSpan(token.Span())

				if isLastToken {
					return nil, newParserError(ErrInvalidSyntaxAndNoRightExpr, token.Span())
				}

				followingToken := input[tokenIdx+1]
				switch followingToken.Kind() {
				case lexer.TokenKindKeywordAND, lexer.TokenKindKeywordOR:
					return nil, newParserError(ErrInvalidSyntaxAndInvalidFollowing, followingToken.Span())
				}
				continue
			}

			node := currentNode.CloneDetached()
			node.WithParent(currentNode)
			node.WithSpan(token.Span())

			currentNode.ClearContent()
			currentNode.WithKind(nodeKindAND)

			if err := currentNode.Insert(node); err != nil {
				return nil, newParserError(fmt.Errorf("%w: %w", ErrInsertChildToNode, err), token.Span())
			}
		}
	}

	if currentNestLevel != 0 {
		lastSpan := input[len(input)-1].Span().End
		return nil, newParserError(ErrInvalidSyntaxUnbalancedParentheses, span.NewSpan(lastSpan, lastSpan))
	}

	expr, err := rootNode.ToExpression()
	if err != nil {
		return nil, err
	}

	return expr, nil
}
