package motor

import "strings"

// Chunk represents a single log write.
type Chunk struct {
	// A formatted message.
	Content string

	// True if a debug write.
	Debug bool

	// A list of inline objects.
	Obj []Object
}

// Object is a structured log element.
//
// They are extracted from named format string operands
// supported by out printf implementation.
//
// See: Brr.Printf()
type Object struct {
	K string
	V interface{}
}

// Endl determines if the chunk ends in a newline.
func (b Chunk) Endl() bool {
	return strings.HasSuffix(b.Content, "\n")
}
