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
}

// Of measures function run time.
func Of(f func()) T {
	t0 := Start()
	f()
	return t0.Stop()
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

// Data returns a list of measurements.
func (b B) Data() Times {
	return b.tf
}

func (b *B) addframe(t T) {
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

// Start constructs a new time measurement frame.
func (b *B) Start() T {
	return T{begin: time.Now().UnixNano()}
}

// Stop
func (b *B) Stop(frame T) T {
	end := time.Now().UnixNano()
	frame.Duration = time.Duration(end - frame.begin)
	b.addframe(frame)
	return frame
}

// Push adds time frames to the benchmark ring.
func (b *B) Push(times ...T) {
	for _, t := range times {
		b.addframe(t)
	}
}

func (b B) Format(state fmt.State, verb rune) {
	t := b.Data()

	// prec, ok := state.Precision()
	// fmt.Println(prec, ok)

	fmt.Fprintf(state, "%v (n=%d, stddev=%v)",
		t.Mean(), len(t), t.Stddev())
}

// Subscribe emits a list of measurements whenever the
// whole ring is fulfilled once.
func (b *B) Subscribe() chan Times {
	if !b.ring {
		panic("speed: can only subscribe for rings")
	}
	if b.done != nil {
		panic("speed: ring is already subscribed to")
	}
	b.done = make(chan struct{})

	ts := make(chan Times)
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
