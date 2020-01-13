//go:generate stringer -type=Keyword
package keyword

type Keyword int

const (
	UNDEFINED Keyword = iota

	// Structural Characters
	BeginArray
	BeginObject
	EndArray
	EndObject
	NameSeparator
	ValueSeparator

	LiteralFalse
	LiteralNull
	LiteralTrue

	Quote
	String
	Number
	Minus
	Plus

	EOF
)
