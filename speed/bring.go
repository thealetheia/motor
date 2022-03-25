package speed

import (
	"fmt"
	"time"
)

// B is a benchmark ring.
//
// Non-zero length B's function as a ring, effectively
// rewriting old frames with the new ones. Zero length
// B's grow normally.
type B struct {
	tf   []T
	ring bool
	pos  int // ring position

	dk, dt stats
}

// Of measures function run time.
func Of(f func()) T {
	t := Now()
	f()
	return t()
}

// Many constructs a benchmark.
func Many(ringSize ...int) B {
	var (
		b    B
		size int = 8
	)
	if ringSize != nil {
		b.ring = true
		size = ringSize[0]
	}
	b.tf = make([]T, 0, size)
	return b
}

// Avg corresponds to the arithmetic average.
func (b B) Avg() (t time.Duration, k float64) {
	return b.Avgt(), b.Avgk()
}

// Avgt returns the average duration of all measurements.
func (b B) Avgt() time.Duration {
	return time.Duration(b.dt.avg)
}

// Avgk returns the average of stored values in the ring.
func (b B) Avgk() float64 {
	return b.dk.avg
}

// Std corresponds to the standard deviation.
func (b B) Std() (t time.Duration, k float64) {
	return b.Stdt(), b.Stdk()
}

// Avgt returns the standard deviation of the durations.
func (b B) Stdt() time.Duration {
	return time.Duration(b.dt.std())
}

// Avgt returns the standard deviation of the K-values.
func (b B) Stdk() float64 {
	return b.dk.std()
}

// Len is the number of items in the underlying slice.
func (b B) Len() int {
	return len(b.tf)
}

// Len is the capacity of the underlying slice.
func (b B) Cap() int {
	return cap(b.tf)
}

// Unordered returns the list of time frames as-is.
func (b B) Unordered() []T {
	return b.tf
}

// Ordered returns the list of frames adjusted for ring position.
func (b B) Ordered() []T {
	if !b.ring || b.pos == 0 {
		return b.tf
	}
	return append(b.tf[b.pos:], b.tf[:b.pos]...)
}

// Adds pushes a time frame to the benchmark ring.
func (b *B) Add(t T, value ...float64) {
	switch len(value) {
	case 0:
	case 1:
		t.K = value[0]
	default:
		panic("speed: too many arguments")
	}
	b.dk.update(t.K)
	b.dt.update(float64(t.Duration))

	if b.pos == len(b.tf) {
		if b.ring && b.pos == cap(b.tf) {
			b.pos = 0
			if len(b.tf) == 0 {
				b.tf = append(b.tf, t)
			}
		} else {
			b.tf = append(b.tf, t)
		}
	}
	b.tf[b.pos] = t
	b.pos++
}

func (b B) Format(state fmt.State, verb rune) {
	avgt, stdt := b.Avgt(), b.Stdt()
	avgk, stdk := b.Avgk(), b.Stdk()

	if avgk == 0.0 {
		fmt.Fprintf(state,
			"%v (n=%d, stdt=%v)",
			avgt, len(b.tf), stdt)
	} else {
		fmt.Fprintf(state,
			"%v/%v (n=%d, stdt=%v, stdk=%g)",
			avgt, avgk, len(b.tf), stdt, stdk)
	}
}
