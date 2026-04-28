package expression

import (
	"testing"

	"github.com/efremenkovan/hlup/span"
	"github.com/stretchr/testify/require"
)

type tExpr struct{}

func (t *tExpr) MatchContains(tokenSequence TokenStream) (bool, []span.Span) {
	return true, []span.Span{}
}

func (t *tExpr) MatchEqual(tokenSequence TokenStream) (bool, []span.Span) {
	return true, []span.Span{}
}

type fExpr struct{}

func (f *fExpr) MatchContains(tokenSequence TokenStream) (bool, []span.Span) {
	return false, []span.Span{}
}

func (f *fExpr) MatchEqual(tokenSequence TokenStream) (bool, []span.Span) {
	return false, []span.Span{}
}

var (
	te Expression = &tExpr{}
	fe Expression = &fExpr{}
)

func Test_And(t *testing.T) {
	seq := ts("any")

	t.Run("match", func(t *testing.T) {
		t.Run("should be true", func(t *testing.T) {
			and := AndExpression{
				Left:  te,
				Right: te,
			}
			match, _ := and.MatchContains(seq)
			require.True(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			and := AndExpression{
				Left:  fe,
				Right: te,
			}
			match, _ := and.MatchContains(seq)
			require.False(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			and := AndExpression{
				Left:  te,
				Right: fe,
			}
			match, _ := and.MatchContains(seq)
			require.False(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			and := AndExpression{
				Left:  fe,
				Right: fe,
			}
			match, _ := and.MatchContains(seq)
			require.False(t, match)
		})
	})

	t.Run("exact match", func(t *testing.T) {
		t.Run("should be true", func(t *testing.T) {
			and := AndExpression{
				Left:  te,
				Right: te,
			}
			match, _ := and.MatchEqual(seq)
			require.True(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			and := AndExpression{
				Left:  fe,
				Right: te,
			}
			match, _ := and.MatchEqual(seq)
			require.False(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			and := AndExpression{
				Left:  te,
				Right: fe,
			}
			match, _ := and.MatchEqual(seq)
			require.False(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			and := AndExpression{
				Left:  fe,
				Right: fe,
			}
			match, _ := and.MatchEqual(seq)
			require.False(t, match)
		})
	})
}

func Test_Or(t *testing.T) {
	seq := ts("any")

	t.Run("match", func(t *testing.T) {
		t.Run("should be true", func(t *testing.T) {
			or := OrExpression{
				Left:  te,
				Right: te,
			}
			match, _ := or.MatchContains(seq)
			require.True(t, match)
		})

		t.Run("should be true", func(t *testing.T) {
			or := OrExpression{
				Left:  fe,
				Right: te,
			}
			match, _ := or.MatchContains(seq)
			require.True(t, match)
		})

		t.Run("should be true", func(t *testing.T) {
			or := OrExpression{
				Left:  te,
				Right: fe,
			}
			match, _ := or.MatchContains(seq)
			require.True(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			or := OrExpression{
				Left:  fe,
				Right: fe,
			}
			match, _ := or.MatchContains(seq)
			require.False(t, match)
		})
	})

	t.Run("exact match", func(t *testing.T) {
		t.Run("should be true", func(t *testing.T) {
			or := OrExpression{
				Left:  te,
				Right: te,
			}
			match, _ := or.MatchEqual(seq)
			require.True(t, match)
		})

		t.Run("should be true", func(t *testing.T) {
			or := OrExpression{
				Left:  fe,
				Right: te,
			}
			match, _ := or.MatchEqual(seq)
			require.True(t, match)
		})

		t.Run("should be true", func(t *testing.T) {
			or := OrExpression{
				Left:  te,
				Right: fe,
			}
			match, _ := or.MatchEqual(seq)
			require.True(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			or := OrExpression{
				Left:  fe,
				Right: fe,
			}
			match, _ := or.MatchEqual(seq)
			require.False(t, match)
		})
	})
}

func Test_Not(t *testing.T) {
	seq := ts("any")

	t.Run("match", func(t *testing.T) {
		t.Run("should be true", func(t *testing.T) {
			ne := NotExpression{
				Expr: fe,
			}
			match, _ := ne.MatchContains(seq)
			require.True(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			ne := NotExpression{
				Expr: te,
			}
			match, _ := ne.MatchContains(seq)
			require.False(t, match)
		})
	})

	t.Run("exact match", func(t *testing.T) {
		t.Run("should be true", func(t *testing.T) {
			ne := NotExpression{
				Expr: fe,
			}
			match, _ := ne.MatchEqual(seq)
			require.True(t, match)
		})

		t.Run("should be false", func(t *testing.T) {
			ne := NotExpression{
				Expr: te,
			}
			match, _ := ne.MatchEqual(seq)
			require.False(t, match)
		})
	})
}
