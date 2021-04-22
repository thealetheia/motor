package speed

import (
	"fmt"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
)

// T is a time frame.
type T struct {
	time.Duration
	// N is meant to represent the frame.
	N float64

	begin int64
}

// Now makes a new time measurement starting now.
//
//     t1 := speed.Now()
//     <-time.After(100*time.Millisecond)
//     debug(t1()) // 101.383ms
//
func Now() func(n ...float64) T {
	t := T{begin: time.Now().UnixNano()}
	return func(n ...float64) T {
		if n != nil {
			t.N = n[0]
		}
		end := time.Now().UnixNano()
		t.Duration = time.Duration(end - t.begin)
		return t
	}
}

func (t T) Begin() time.Time {
	return time.Unix(0, t.begin)
}

func (t T) End() time.Time {
	return time.Unix(0, t.begin+int64(t.Duration))
}

func (t T) Format(state fmt.State, verb rune) {
	dur := int64(t.Duration)

	if state.Flag('+') {
		fmt.Fprint(state, t.begin, "-")
		if dur != 0 {
			fmt.Fprint(state, t.begin+dur)
			goto suffix
		}
	}
	if dur == 0 {
		fmt.Fprint(state, "???")
	} else {
		fmt.Fprintf(state, "%v", t.Round(time.Microsecond))
	}

suffix:
	if t.N != 0 {
		n := strconv.FormatFloat(t.N, 'f', -1, 64)
		fmt.Fprintf(state, "(n=%s)", n)
	}
}

// Frames i an ordered list of time measurements.
type Frames []T

func (tf Frames) Mean() time.Duration {
	t := make([]float64, len(tf))
	for i := range t {
		t[i] = float64(tf[i].Duration)
	}
	mean, _ := stats.Mean(t)
	return time.Duration(mean).Round(time.Microsecond)
}

func (tf Frames) Stddev() time.Duration {
	t := make([]float64, len(tf))
	for i := range t {
		t[i] = float64(tf[i].Duration)
	}
	sdev, _ := stats.StdDevP(t)
	return time.Duration(sdev).Round(time.Microsecond)
}
