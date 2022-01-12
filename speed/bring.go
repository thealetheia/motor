package speed

import (
	"fmt"
	"time"
)

// B is a benchmark: repeated time measurement.
//
// Non-zero length B's function as a ring, effectively
// rewriting old frames with the new ones. Zero length
// B's grow normally.
type B struct {
	tf   []T
	ring bool
	pos  int // ring position
	done chan struct{}

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
func (b B) Avgt() time.Duration {
	return time.Duration(b.dt.avg)
}
func (b B) Avgk() float64 {
	return b.dk.avg
}

// Std corresponds to the standard deviation.
func (b B) Std() (t time.Duration, k float64) {
	return b.Stdt(), b.Stdk()
}
func (b B) Stdt() time.Duration {
	return time.Duration(b.dt.std())
}
func (b B) Stdk() float64 {
	return b.dk.std()
}

// All returns a list of measurements.
func (b B) All() []T {
	return b.tf
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

	if b.ring && b.done != nil && b.pos == cap(b.tf) {
		b.done <- struct{}{}
	}
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

// Subscribe emits a list of measurements whenever the
// whole ring is fulfilled once.
func (b *B) Subscribe() chan []T {
	if !b.ring {
		panic("speed: can only subscribe for rings")
	}
	if b.done != nil {
		panic("speed: ring is already subscribed to")
	}
	b.done = make(chan struct{})

	ts := make(chan []T)
	go func() {
		for range b.done {
			t := make([]T, len(b.tf))
			copy(t, b.tf)
			b.tf = b.tf[:]
			b.pos = 0
			ts <- t
		}
		close(ts)
	}()
	return ts
}
