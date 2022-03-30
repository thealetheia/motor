package motor

import (
	"fmt"
	"io"
	"strings"
)

// Chunk represents a single log write to the log buffer.
type Chunk struct {
	// Printf format string, if provided.
	Format string

	Args []interface{}
	// Named printf arguments (tags) list.
	Tags  []Tag
	Flags []Flag

	// True if the chunk is a debug write in itself.
	Debug bool
}

// Autowrite prints (and formats) chunk arguments.
func (c Chunk) Autowrite(w io.Writer) {
	if c.Format == "" {
		fmt.Fprintln(w, c.Args...)
	} else {
		f := c.Format
		if !strings.HasSuffix(f, "\n") {
			f += "\n"
		}
		fmt.Fprintf(w, c.Format, c.Args...)
	}
}

// Flag returns true if chunk contains flag.
func (c Chunk) Flag(f Flag) bool {
	for i := range c.Flags {
		if c.Flags[i] == f {
			return true
		}
	}
	return false
}

// Flag allows for customization of chunks.
type Flag int

const (
	Error Flag = iota
	// Pretty yellow warning thing :-)
	Warn
	Trace
)
