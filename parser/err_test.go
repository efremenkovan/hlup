package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrInvalidSyntax(t *testing.T) {
	invlidSyntaxErrs := []error{
		ErrInvalidSyntax,
		ErrInvalidSyntaxUnbalancedParentheses,
		ErrInvalidSyntaxNotNoFollowingExpr,
		ErrInvalidSyntaxNotInvalidFollowing,
		ErrInvalidSyntaxAndNoLeftExpr,
		ErrInvalidSyntaxAndNoRightExpr,
		ErrInvalidSyntaxAndInvalidFollowing,
		ErrInvalidSyntaxOrNoLeftExpr,
		ErrInvalidSyntaxOrNoRightExpr,
		ErrInvalidSyntaxOrInvalidFollowing,
	}

	for _, err := range invlidSyntaxErrs {
		t.Run("all invalid syntax errors should be checkable by errors.Is(..., ErrInvalidSyntax)", func(t *testing.T) {
			require.True(t, errors.Is(err, ErrInvalidSyntax))
		})
	}
}
