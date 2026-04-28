package lexer

type ctxBit = int8

const (
	// Lexer is now within literal, surrounded by quotes
	ctxInQuote ctxBit = 1 << 0

	// Lexer should escape current cursor rune
	ctxEscapedRune ctxBit = 1 << 1
)

type ctx struct {
	state ctxBit

	// When in `InQuote` context, this rune means InQuote context termination
	terminalQuote rune
}

func newCtx() ctx {
	return ctx{
		state: 0,
	}
}

// Is returns true if context state contains specified bit
func (c *ctx) Is(unit ctxBit) bool {
	return (c.state & unit) != 0
}

// Add mutates context state by extending it with provided bit
func (c *ctx) Add(unit ctxBit) {
	c.state |= unit
}

// Drop mutates context state by removing provided bit
func (c *ctx) Drop(unit ctxBit) {
	c.state &^= unit
}

// Clone creates copy of the context with its state
func (c *ctx) Clone() ctx {
	return ctx{
		state:         c.state,
		terminalQuote: c.terminalQuote,
	}
}
