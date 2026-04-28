package lang

type Lang int8

// List of all supported languages
const (
	LangRU Lang = iota
	LangEN
)

// String returns a string representations of language
func (s Lang) String() string {
	switch s {
	case LangRU:
		return "ru"
	case LangEN:
		return "en"
	}

	return "unknown language"
}
