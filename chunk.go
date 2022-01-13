package motor

import "strings"

// Chunk represents a single log write.
type Chunk struct {
	Fmt  string
	Args []interface{}

	Debug bool
}

// Endl determines if the chunk ends in a newline.
func (b Chunk) Endl() bool {
	return strings.HasSuffix(b.Fmt, "\n")
}
