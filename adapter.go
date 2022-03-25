package motor

import (
	"io"
	"time"
)

// Adapter is how motor interacts with the destination I/O.
//
type Adapter interface {
	// Begin is called whenever a procedure is started.
	Begin(*Brr)

	// Write manages the log transformation.
	//
	// This function transforms the provided chunk into
	// the intermediate batching buffer in charge of the
	// flush valve.
	Write(brr *Brr, chunk Chunk, w io.Writer)

	// End is called just before the final flush.
	End(*Brr)

	// "Real" destination of the log stream.
	//
	// The log writes are flushed here at the most suitable
	// time to minimize memory use and I/O pressure due to
	// irregular writes typical for logging.
	Device() io.Writer

	// Max allowed time-to-flush for any given write.
	AllowedLatency() time.Duration

	// Only let through tagged writes.
	Tagged() bool
}
