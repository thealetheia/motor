package speed

import (
	"fmt"
	"strconv"
	"time"

	"github.com/montanaflynn/stats"
)

const (
	Ns = time.Nanosecond
	Us = time.Microsecond
	Ms = time.Millisecond
)

// T is a time frame.
type T struct {
	time.Duration
	// N is meant to represent the frame.
	N float64

	begin int64
}

// Start makes a new time measurement starting now.
//
//     t0 := speed.Start()
//     <-time.After(100*time.Millisecond)
//     debug(t0.Stop()) // 101.383ms
//
func Start() T {
	return T{begin: time.Now().UnixNano()}
}

func (t T) Stop() T {
	end := time.Now().UnixNano()
	t.Duration = time.Duration(end - t.begin)
	return t
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

// Times i an ordered list of time measurements.
type Times []T

func (ts Times) Mean() time.Duration {
	t := make([]float64, len(ts))
	for i := range t {
		t[i] = float64(ts[i].Duration)
	}
	mean, _ := stats.Mean(t)
	return time.Duration(mean).Round(time.Microsecond)
}

func (ts Times) Stddev() time.Duration {
	t := make([]float64, len(ts))
	for i := range t {
		t[i] = float64(ts[i].Duration)
	}
	sdev, _ := stats.StdDevP(t)
	return time.Duration(sdev).Round(time.Microsecond)
}
