package speed

import (
	"fmt"
	"time"
)

// F is a time frame.
type F struct {
	// A number meant to represent the frame.
	N     float64
	Begin time.Time
	End   time.Time

	// Internal frame id within T.
	id int
}

func (f F) Format(state fmt.State, verb rune) {
	fmt.Printf("verb %s\n", string(verb))
}

// Dt is the ∆t frame duration.
func (f F) Dt() time.Duration {
	return f.End.Sub(f.Begin)
}

// Frames i an ordered list of time measurements.
type Frames []F
