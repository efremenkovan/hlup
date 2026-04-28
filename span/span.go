package span

// Span represents a reference to a range of runes in input string, associated with some computed structure
type Span struct {
	Start int
	End   int
}

func NewSpan(start, end int) Span {
	return Span{Start: start, End: end}
}
