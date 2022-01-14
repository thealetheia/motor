package motor

import (
	"fmt"
	"io"
)

// Chunk represents a single log write to the log buffer.
type Chunk struct {
	// Printf format string, if provided.
	Format string

	// Print arguments, if provided.
	Args  []interface{}
	Tags  []Tag
	Flags []Flag

	// True if the chunk is a debug message.
	Debug bool
}

// Autowrite prints (and formats) the chunk arguments.
func (c Chunk) Autowrite(w io.Writer) {
	if c.Format == "" {
		fmt.Fprintln(w, c.Args...)
	} else {
		fmt.Fprintf(w, c.Format, c.Args...)
	}
}
