package speed

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
)

// F is a time frame.
type F struct {
	time.Duration

	// A number meant to represent the frame.
	N     float64
	Begin time.Time
	End   time.Time

	// Internal frame id within T.
	id int
}

// Now makes a new time measurement starting now.
//
//     t1 := speed.Now()
//     <-time.After(100*time.Millisecond)
//     debug(t1()) // 101.383ms
//
func Now() func(n ...float64) F {
	f := F{Begin: time.Now()}
	return func(n ...float64) F {
		if n != nil {
			f.N = n[0]
		}
		f.End = time.Now()
		f.Duration = f.dt()
		return f
	}
}

func (f F) Format(state fmt.State, verb rune) {
	var b bytes.Buffer
	if state.Flag('+') {
		fmt.Fprint(&b, f.Begin.UnixNano(), "-")
		if !f.End.IsZero() {
			fmt.Fprint(&b, f.End.UnixNano())
			goto format
		}
	}
	if f.End.IsZero() {
		b.WriteString("???")
	} else {
		fmt.Fprintf(&b, "%v", f.dt().Round(time.Microsecond))
	}

format:
	if f.N != 0 {
		n := strconv.FormatFloat(f.N, 'f', -1, 64)
		fmt.Fprintf(&b, "(n=%s)", n)
	}
	b.WriteTo(state)
}

// ∆t
func (f F) dt() time.Duration {
	return f.End.Sub(f.Begin)
}

// Frames i an ordered list of time measurements.
type Frames []F

func (fs Frames) Mean() time.Duration {
	t := make([]float64, len(fs))
	for i := range t {
		t[i] = float64(fs[i].dt())
	}
	mean, _ := stats.Mean(t)
	return time.Duration(mean).Round(time.Microsecond)
}

func (fs Frames) Stddev() time.Duration {
	t := make([]float64, len(fs))
	for i := range t {
		t[i] = float64(fs[i].dt())
	}
	sdev, _ := stats.StdDevP(t)
	return time.Duration(sdev).Round(time.Microsecond)
}
