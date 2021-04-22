package speed

import (
	"fmt"
	"time"
)

// B is a benchmark.
//
// Benchmarks are repeated time measurements.
//
// Non-zero length B's function as a ring, effectively
// rewriting old frames with the new ones. Zero length
// B's grow normally.
type B struct {
	tf   []T
	ring bool
	// ring position
	pos int
}

// Of constructs a benchmark.
func Of(ringSize ...int) *B {
	b := &B{}
	var size int = 16
	if ringSize != nil {
		b.ring = true
		size = ringSize[0]
	}
	b.tf = make([]T, 0, size)
	return b
}

// Frames returns an ordered list of past measurements.
func (b *B) Frames() Frames {
	if b.ring && b.pos < len(b.tf) {
		if b.tf[b.pos].begin > b.tf[b.pos+1].begin {
			return Frames(append(b.tf[b.pos:], b.tf[:b.pos]...))
		}
	}

	return b.tf
}

func (b *B) addframe() (T, int) {
	t := T{begin: time.Now().UnixNano()}
	id := b.pos
	if b.pos == len(b.tf) {
		if b.ring && b.pos == cap(b.tf) {
			b.pos = 0
			id = 0
			if len(b.tf) == 0 {
				b.tf = append(b.tf, t)
			}
		} else {
			b.tf = append(b.tf, t)
		}
	}
	b.tf[b.pos] = t
	b.pos++
	return t, id
}

// Now constructs a new time measurement frame.
func (b *B) Now() func(n ...float64) T {
	t, id := b.addframe()
	return func(n ...float64) T {
		end := time.Now().UnixNano()
		t.Duration = time.Duration(end - t.begin)
		if n != nil {
			t.N = n[0]
		}

		if b.tf[id].begin != t.begin {
			return t
		}
		b.tf[id] = t
		return t
	}
}

// Push adds a time frame to the benchmark ring.
func (b *B) Push(frame T) {
	_, id := b.addframe()
	b.tf[id] = frame
}

func (b *B) Format(state fmt.State, verb rune) {
	tf := b.Frames()
	fmt.Fprintf(state, "%v (n=%d, stddev=%v)",
		tf.Mean(), len(tf), tf.Stddev())
}

// Done returns a list of frames whenever the whole ring
// is fulfilled once.
// func (t *T) Done() chan Frames {
// 	if !b.ring {
// 		panic("speed: can only wait for rings")
// 	}
// 	b.done = make(chan struct{})
// 	go func() {
// 		for _, f := range b.f {
//			logic
// 		}
// 		close(done)
// 	}()
// 	return done
// }
