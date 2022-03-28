package motor

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"aletheia.icu/motor/speed"
	"github.com/fatih/color"
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
	//
	// Returns negative if unlimited.
	MaxLatency() time.Duration

	// Indicates whether if the adapter collects untagged
	// write chunks.
	//
	// It might not, i.e. when collecting telemetry only.
	TaggerStagger() bool
}

// Simple is a basic adapter with terminal formatting.
type Simple struct {
}

func (adp *Simple) Begin(brr *Brr) {
	bold := color.New(color.Bold)
	bold.Print("BEGIN")
	fmt.Printf(" %s %s\n",
		brr.Namekind, brr.Id)
}

func (adp *Simple) Write(brr *Brr, c Chunk, w io.Writer) {
	fmt.Fprintf(w, "[%s+%v] ", brr.Id,
		speed.Trunc(time.Now().Sub(brr.T0), 4))

	if c.Format == "" {
		fmt.Fprintln(w, c.Args...)
		return
	}

	fmt.Fprintf(w, c.Format, c.Args...)
	if c.Tags != nil {
		tagp := make([]string, len(c.Tags))
		for i := range tagp {
			tagp[i] = c.Tags[i].Label
		}
		fmt.Fprintf(w, " [%s]", strings.Join(tagp, ","))
	}

	if c.Debug {
		fmt.Fprint(w, "*")
	}
	fmt.Fprintln(w)
}

func (adp *Simple) End(brr *Brr) {
	bold := color.New(color.Bold)
	bold.Print("END")
	fmt.Printf(" %s %s AFTER %v\n",
		brr.Namekind, brr.Id, time.Now().Sub(brr.T0))
}

func (adp *Simple) Device() io.Writer {
	return os.Stdout
}

func (adp *Simple) MaxLatency() time.Duration {
	return -time.Second
}

func (adp *Simple) TaggerStagger() bool {
	return false
}
