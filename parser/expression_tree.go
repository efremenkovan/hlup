package parser

import (
	"github.com/efremenkovan/hlup/expression"
	"github.com/efremenkovan/hlup/lexer"
	"github.com/efremenkovan/hlup/span"
)

type nodeKind int8

const (
	nodeKindUnknown nodeKind = iota
	nodeKindOR
	nodeKindAND
	nodeKindNOT
	nodeKindLeaf
)

type node struct {
	parent iNode
	kind   nodeKind
	span   span.Span

	lNode iNode
	rNode iNode
}

// iNode represents an expression tree node. It can be an expression node (and/or/not) or it can be a leaf node (plain expression token)
type iNode interface {
	Left() iNode
	Right() iNode
	Kind() nodeKind
	Span() span.Span

	// Parent returns a pointer to a parent node
	// If called on the root node returns nil
	Parent() iNode

	// ToExpression returns an expression, represented by this node
	ToExpression() (expression.Expression, error)

	// IsRoot returns true if node is root node
	IsRoot() bool

	// IsEmpty returns true if both the Left and the Right child nodes are unspecified
	IsEmpty() bool

	// IsFull returns true if both the Left and the Right child nodes are specified
	IsFull() bool

	// WithKind mutates the node attaching a specified kind to it
	WithKind(kind nodeKind)

	// WithSpan mutates the node attaching a specified span to it
	WithSpan(span span.Span)

	// WithParent mutates the node attaching a specified parent to it
	WithParent(parent iNode)

	// CloneDetached creates and returns a new node with no parent and contents identical to this node
	CloneDetached() iNode

	// ClearContent mutates the node resetting its kind and child nodes
	ClearContent()

	// Insert mutates the node inserting provided node as child node.
	// If both the Left and the Right branches are already specified returns an error.
	Insert(node iNode) error
}

func newNode() iNode {
	return &node{
		parent: nil,
		kind:   nodeKindUnknown,

		// span of the keyword
		span: span.Span{},

		lNode: nil,
		rNode: nil,
	}
}

func (n *node) Span() span.Span {
	switch n.kind {
	case nodeKindNOT:
		return span.NewSpan(n.span.Start, n.lNode.Span().End)
	case nodeKindAND, nodeKindOR:
		return span.NewSpan(n.lNode.Span().Start, n.rNode.Span().End)
	}

	return n.span
}

func (n *node) WithSpan(span span.Span) {
	n.span = span
}

func (n *node) IsRoot() bool {
	return n.parent == nil
}

func (n *node) Parent() iNode {
	return n.parent
}

func (n *node) IsEmpty() bool {
	return n.lNode == nil && n.rNode == nil
}

func (n *node) IsFull() bool {
	return (n.lNode != nil) && (n.rNode != nil)
}

func (n *node) WithKind(kind nodeKind) {
	n.kind = kind
}

func (n *node) WithParent(parent iNode) {
	n.parent = parent
}

func (n *node) Kind() nodeKind {
	return n.kind
}

func (n *node) Left() iNode {
	return n.lNode
}

func (n *node) Right() iNode {
	return n.rNode
}

func (n *node) CloneDetached() iNode {
	node := newNode()
	node.WithKind(n.kind)
	if n.lNode != nil {
		n.lNode.WithParent(node)
		_ = node.Insert(n.lNode)
	}

	if n.rNode != nil {
		n.rNode.WithParent(node)
		_ = node.Insert(n.rNode)
	}

	return node
}

func (n *node) ClearContent() {
	n.kind = nodeKindUnknown
	n.lNode = nil
	n.rNode = nil
}

func (n *node) Insert(node iNode) error {
	if n.lNode == nil {
		n.lNode = node
		return nil
	}

	if n.rNode == nil {
		n.rNode = node
		return nil
	}

	return ErrTreeInsertIntoFullNode
}

func (n node) ToExpression() (expression.Expression, error) {
	switch n.kind {
	case nodeKindNOT:
		if n.lNode == nil {
			return nil, ErrToExpressionNoLeftNode
		}

		expr, err := n.lNode.ToExpression()
		if err != nil {
			return nil, err
		}

		return expression.NotExpression{
			Expr: expr,
		}, nil
	case nodeKindAND:
		if n.lNode == nil {
			return nil, ErrToExpressionNoLeftNode
		}
		if n.rNode == nil {
			return nil, ErrToExpressionNoRightNode
		}

		lExpr, err := n.lNode.ToExpression()
		if err != nil {
			return nil, err
		}

		rExpr, err := n.rNode.ToExpression()
		if err != nil {
			return nil, err
		}

		return expression.AndExpression{
			Left:  lExpr,
			Right: rExpr,
		}, nil
	case nodeKindOR:
		if n.lNode == nil {
			return nil, ErrToExpressionNoLeftNode
		}
		if n.rNode == nil {
			return nil, ErrToExpressionNoRightNode
		}

		lExpr, err := n.lNode.ToExpression()
		if err != nil {
			return nil, err
		}

		rExpr, err := n.rNode.ToExpression()
		if err != nil {
			return nil, err
		}

		return expression.OrExpression{
			Left:  lExpr,
			Right: rExpr,
		}, nil

		// This case handles root level expression wrapped in parentheses
	case nodeKindUnknown:
		// (...) can only result with single left child node. Other variants are unknown behavior
		if (n.lNode == nil && n.rNode == nil) || (n.lNode != nil && n.rNode != nil) {
			return nil, ErrToExpressionUnknownNodeKind
		}

		switch n.lNode.Kind() {
		case nodeKindUnknown, nodeKindLeaf:
			return nil, ErrToExpressionUnknownNodeKind
		}

		return n.lNode.ToExpression()
	}

	return nil, ErrToExpressionUnknownNodeKind
}

type leaf struct {
	tokens lexer.TokenStream
	parent iNode
}

func (l *leaf) ToExpression() (expression.Expression, error) {
	res := make(expression.TokenStream, len(l.tokens))
	for i, token := range l.tokens {
		res[i] = expression.Token{Value: token.String(), Span: token.Span()}
	}

	return res, nil
}

func newLeaf(token lexer.Token) leaf {
	return leaf{
		parent: nil,
		tokens: lexer.TokenStream{token},
	}
}

func emptyLeaf() leaf {
	return leaf{
		parent: nil,
		tokens: lexer.TokenStream{},
	}
}

func (l *leaf) Span() span.Span {
	if len(l.tokens) == 0 {
		return span.Span{}
	}
	return span.Span{
		Start: l.tokens[0].Span().Start,
		End:   l.tokens[len(l.tokens)-1].Span().End,
	}
}

func (l *leaf) WithSpan(span span.Span) {
	// noop
}

func (l *leaf) Left() iNode {
	return nil
}

func (l *leaf) Right() iNode {
	return nil
}

func (l *leaf) IsRoot() bool {
	return false
}

func (l *leaf) IsEmpty() bool {
	return len(l.tokens) == 0
}

func (l *leaf) IsFull() bool {
	return false
}

func (l *leaf) WithKind(kind nodeKind) {
	// noop
}

func (l *leaf) WithParent(parent iNode) {
	l.parent = parent
}

func (l *leaf) CloneDetached() iNode {
	return &leaf{
		tokens: l.tokens,
	}
}

func (l *leaf) Parent() iNode {
	return l.parent
}

func (l *leaf) Kind() nodeKind {
	return nodeKindLeaf
}

func (l *leaf) ClearContent() {
	l.tokens = lexer.TokenStream{}
}

func (l *leaf) Insert(node iNode) error {
	return ErrTreeInsertIntoLeaf
}

func (l *leaf) AppendContent(token lexer.Token) {
	l.tokens = append(l.tokens, token)
}
