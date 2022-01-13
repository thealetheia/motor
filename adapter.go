package motor

import (
	"io"
)

// Adapter is how motor interacts with destination I/O.
//
type Adapter interface {
	// Begin is called whenever a procedure has started.
	Begin(*Brr) error

	// Write manages the log transformation.
	//
	// This function transforms the provided chunk into
	// the intermediate batching buffer in charge of the
	// flush valve.
	Write(brr *Brr, chunk Chunk, w io.Writer) error

	// End is called after the final flush.
	End(*Brr) error

	// "Real" destination of the log stream.
	//
	// All writes will be flushed here eventually.
	Device() io.Writer
}
