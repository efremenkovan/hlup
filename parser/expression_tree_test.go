package parser

import (
	"testing"

	"github.com/efremenkovan/hlup/lexer"
	"github.com/efremenkovan/hlup/span"
	"github.com/stretchr/testify/require"
)

var (
	dummyNode = newNode()
	dummyLeaf = newLeaf(lexer.NewLiteralToken("", span.NewSpan(0, 0)))
)

func Test_Node(t *testing.T) {
	t.Run("should consider a node without a parent as a root node", func(t *testing.T) {
		node := newNode()
		require.True(t, node.IsRoot())
	})

	t.Run("consider a node full", func(t *testing.T) {
		type test struct {
			name string
			node node
		}

		tests := []test{
			{name: "contains two nodes", node: node{lNode: dummyNode, rNode: dummyNode}},
			{name: "contains two leafs", node: node{lNode: &dummyLeaf, rNode: &dummyLeaf}},
			{name: "contains left leafs right node", node: node{lNode: &dummyLeaf, rNode: dummyNode}},
			{name: "contains left node right leaf", node: node{lNode: dummyNode, rNode: &dummyLeaf}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.True(t, tt.node.IsFull())
			})
		}
	})

	t.Run("do not consider a node full", func(t *testing.T) {
		type test struct {
			name string
			node node
		}

		tests := []test{
			{name: "contains one node", node: node{lNode: dummyNode}},
			{name: "contains one leaf", node: node{lNode: &dummyLeaf}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.False(t, tt.node.IsFull())
			})
		}
	})

	t.Run("consider a node empty", func(t *testing.T) {
		node := newNode()
		require.True(t, node.IsEmpty())
	})

	t.Run("do not consider a node empty", func(t *testing.T) {
		type test struct {
			name string
			node node
		}

		tests := []test{
			{name: "contains one node", node: node{lNode: dummyNode}},
			{name: "contains one leaf", node: node{lNode: &dummyLeaf}},
			{name: "contains two nodes", node: node{lNode: dummyNode, rNode: dummyNode}},
			{name: "contains two leafs", node: node{lNode: &dummyLeaf, rNode: &dummyLeaf}},
			{name: "contains left leafs right node", node: node{lNode: &dummyLeaf, rNode: dummyNode}},
			{name: "contains left node right leaf", node: node{lNode: dummyNode, rNode: &dummyLeaf}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.False(t, tt.node.IsEmpty())
			})
		}
	})

	t.Run("clone detached should have same content, other pointer address and no parent", func(t *testing.T) {
		left := newLeaf(lexer.NewLiteralToken("one", span.NewSpan(0, 2)))
		right := newLeaf(lexer.NewLiteralToken("two", span.NewSpan(4, 6)))
		n := node{
			kind:   nodeKindAND,
			lNode:  &left,
			rNode:  &right,
			parent: dummyNode,
		}

		newNode := n.CloneDetached()
		require.Equal(t, &left, newNode.Left())
		require.Equal(t, &right, newNode.Right())
		require.Nil(t, newNode.Parent())
	})

	t.Run("clone detached should change left and right nodes' parent to new", func(t *testing.T) {
		left := newNode()
		right := newNode()
		n := node{
			kind:   nodeKindAND,
			lNode:  left,
			rNode:  right,
			parent: dummyNode,
		}

		newNode := n.CloneDetached()
		require.Equal(t, left, newNode.Left())
		require.Equal(t, left.Parent(), newNode)
		require.Equal(t, right, newNode.Right())
		require.Equal(t, right.Parent(), newNode)
		require.Nil(t, newNode.Parent())
	})

	t.Run("clear content should clear node's content leaving the parent unchanged", func(t *testing.T) {
		left := newNode()
		right := newNode()
		n := node{
			kind:   nodeKindAND,
			lNode:  left,
			rNode:  right,
			parent: dummyNode,
		}

		n.ClearContent()
		require.Equal(t, n.kind, nodeKindUnknown)
		require.Nil(t, n.lNode)
		require.Nil(t, n.rNode)
		require.Equal(t, dummyNode, n.parent)
	})
}
